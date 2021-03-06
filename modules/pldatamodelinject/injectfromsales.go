package main

import (
	"eaciit/gdrj/model"
	"eaciit/gdrj/modules"
	"os"

	"sync"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm/v1"
	"github.com/eaciit/toolkit"
)

var conn dbox.IConnection
var count int
var wg *sync.WaitGroup
var mtx *sync.Mutex

func setinitialconnection() {
	var err error
	conn, err = modules.GetDboxIConnection("db_godrej")

	if err != nil {
		toolkit.Println("Initial connection found : ", err)
		os.Exit(1)
	}

	err = gdrj.SetDb(conn)
	if err != nil {
		toolkit.Println("Initial connection found : ", err)
		os.Exit(1)
	}
}

var (
	ledgers = toolkit.M{}
	/*
	   shs = toolkit.M{}
	   pcs = toolkit.M{}
	   ccs = toolkit.M{}
	   prods = toolkit.M{}
	   custs = toolkit.M{}
	   vdistskus = toolkit.M{}
	*/
)

func getCursor(obj orm.IModel) dbox.ICursor {
	c, e := gdrj.Find(obj, nil, nil)
	if e != nil {
		return nil
	}
	return c
}

func prepMaster() {
	var e error
	/*
	   pc:=new(gdrj.ProfitCenter)
	   cc:=new(gdrj.CostCenter)
	   prod:=new(gdrj.Product)

	   cpc := getCursor(pc)
	   defer cpc.Close()
	   for e=cpc.Fetch(pc,1,false);e==nil;{
	       pcs.Set(pc.ID,pc)
	       pc =new(gdrj.ProfitCenter)
	       e=cpc.Fetch(pc,1,false)
	   }

	   ccc:=getCursor(cc)
	   defer ccc.Close()
	   for e=ccc.Fetch(cc,1,false);e==nil;{
	       ccs.Set(cc.ID,cc)
	       cc = new(gdrj.CostCenter)
	       e=ccc.Fetch(cc,1,false)
	   }

	   cprod:=getCursor(prod)
	   defer cprod.Close()
	   for e=cprod.Fetch(prod,1,false);e==nil;{
	       prods.Set(prod.ID,prod)
	       prod=new(gdrj.Product)
	       e=cprod.Fetch(prod,1,false)
	   }

	   cust := new(gdrj.Customer)
	   ccust:=getCursor(cust)
	   defer ccust.Close()
	   for e=ccust.Fetch(cust,1,false);e==nil;{
	       custs.Set(cust.ID,cust)
	       cust=new(gdrj.Customer)
	       e=ccust.Fetch(cust,1,false)
	   }

	   sku:=new(gdrj.MappingInventory)
	   cskus:=getCursor(sku)
	   defer cskus.Close()
	   for e=cskus.Fetch(sku,1,false);e==nil;{
	       vdistskus.Set(sku.SKUID_VDIST,sku.ID)
	       sku=new(gdrj.MappingInventory)
	       e=cskus.Fetch(sku,1,false)
	   }

	   sh := new(gdrj.SalesHeader)
	   cshs:=getCursor(sh)
	   defer cshs.Close()
	   for e=cshs.Fetch(sh,1,false);e==nil;{
	       //sh.SalesGrossAmount=0
	       //sh.SalesNetAmount=0
	       //sh.SalesLine=0
	       shs.Set(sh.ID,sh)
	       sh = new(gdrj.SalesHeader)
	       e=cshs.Fetch(sh,1,false)
	   }
	*/

	ledger := new(gdrj.LedgerMaster)
	cledger := getCursor(ledger)
	defer cledger.Close()
	for e = cledger.Fetch(ledger, 1, false); e == nil; {
		ledgers.Set(ledger.ID, ledger)
		ledger = new(gdrj.LedgerMaster)
		e = cledger.Fetch(ledger, 1, false)
	}
}

func main() {
	setinitialconnection()
	defer gdrj.CloseDb()

	toolkit.Println("Reading Master")
	prepMaster()

	toolkit.Printfn("Delete existing")
	conn.NewQuery().From(new(gdrj.PLDataModel).TableName()).Where(dbox.Eq("source", "SalesVDist")).Delete().Exec(nil)

	toolkit.Println("START...")
	crx, err := gdrj.Find(new(gdrj.SalesTrx),
		dbox.Eq("pcvalid", true),
		//toolkit.M{}.Set("take",1200))
		nil)
	if err != nil {
		toolkit.Println("Error Found : ", err.Error())
		os.Exit(1)
	}
	defer crx.Close()

	count := crx.Count()
	i := 0
	t0 := time.Now()
	models := map[string]*gdrj.PLDataModel{}
	for {
		//st := new(gdrj.SalesTrx)
		sts := []gdrj.SalesTrx{}
		err = crx.Fetch(&sts, 1000, false)
		if err != nil {
			toolkit.Printfn("Exit loop: %s", err.Error())
			break
		}

		if len(sts) > 0 {
			for _, st := range sts {

				i++
				if !st.HeaderValid {
					continue
				}

				ledgeraccount := ""
				if st.SalesQty > 0 {
					ledgeraccount = "70000000"
				} else {
					ledgeraccount = "70000302"
				}

				ledg := ledgers.Get(ledgeraccount).(*gdrj.LedgerMaster)
				idsales := toolkit.Sprintf("%d_%d_%s_%s_%s_%s_%s_%s",
					st.Date.Year, st.Date.Month,
					"ID11",
					ledg.PLCode, st.OutletID, st.SKUID, st.PC.ID, "")
				mdl, bmodel := models[idsales]
				if !bmodel {
					mdl = new(gdrj.PLDataModel)
					mdl.CompanyCode = "ID11"
					mdl.PC = st.PC
					mdl.PCID = st.PC.ID
					mdl.OutletID = st.OutletID
					mdl.SKUID = st.SKUID
					mdl.Customer = st.Customer
					mdl.Product = st.Product
					//mdl.LedgerAccount = ledgeraccount
					gdate := gdrj.NewDate(st.Date.Year(), int(st.Date.Month()), 1)
					mdl.Date = gdate
					mdl.Year = gdate.Year
					mdl.Month = gdate.Month
					ledg := ledgers.Get(ledgeraccount).(*gdrj.LedgerMaster)
					mdl.PLCode = ledg.PLCode
					mdl.PLOrder = ledg.OrderIndex
					mdl.PLGroup1 = ledg.H1
					mdl.PLGroup2 = ledg.H2
					mdl.PLGroup3 = ledg.H3
					mdl.Source = "SalesVDist"
					mdl.Value1 = 0
					mdl.Value2 = 0
					mdl.Value3 = 0
					models[idsales] = mdl
				}
				mdl.Value1 += st.GrossAmount
				mdl.Value3 += st.SalesQty
				mdl.Value2 = mdl.Value1 / mdl.Value3
				//gdrj.Save(mdl)

				if st.DiscountAmount != 0 {
					ledgeraccount = "75053730"
					ledg = ledgers.Get(ledgeraccount).(*gdrj.LedgerMaster)
					iddiscount := toolkit.Sprintf("%d_%d_%s_%s_%s_%s_%s_%s",
						st.Date.Year, st.Date.Month,
						"ID11",
						ledg.PLCode, st.OutletID, st.SKUID, st.PC.ID, "")
					mdisc, bdisc := models[iddiscount]
					if !bdisc {
						mdisc = new(gdrj.PLDataModel)
						*mdisc = *mdl
						mdisc.PLCode = ledg.PLCode
						mdisc.PLOrder = ledg.OrderIndex
						mdisc.PLGroup1 = ledg.H1
						mdisc.PLGroup2 = ledg.H2
						mdisc.PLGroup3 = ledg.H3
						mdisc.Value1 = 0
						mdisc.Value2 = 0
						mdisc.Value3 = 0
						mdisc.Source = "SalesVDist"
						models[iddiscount] = mdisc
					}
					mdisc.Value1 += st.DiscountAmount
					mdisc.Value3 += st.SalesQty
					mdisc.Value2 = mdisc.Value1 / mdisc.Value3
				}

				id_netsales := toolkit.Sprintf("%d_%d_%s_%s_%s_%s_%s_%s",
					st.Date.Year, st.Date.Month,
					"ID11",
					"PL8A", st.OutletID, st.SKUID, st.PC.ID, "")
				mns, bdisc := models[id_netsales]
				if !bdisc {
					mns = new(gdrj.PLDataModel)
					*mns = *mdl
					mns.PLCode = "PL08A"
					mns.PLOrder = "PL009"
					mns.PLGroup1 = "Net Sales"
					mns.PLGroup2 = "Net Sales"
					mns.PLGroup3 = "Net Sales"
					mns.Source = "SalesVDist"
					mns.Value1 = 0
					mns.Value2 = 0
					mns.Value3 = 0
					models[id_netsales] = mns
				}
				mns.Value1 += st.GrossAmount + st.DiscountAmount
				mns.Value3 += st.SalesQty
				mns.Value2 = mns.Value1 / mns.Value3
			}
		}

		toolkit.Printfn("Calc %d of %din %s", i, count, 
			time.Since(t0).String())

		if len(sts)<1000{
			break
		}
	}

	count = len(models)
	i = 0
	toolkit.Printfn("Saving data %d records", count)

	jobs := make(chan *gdrj.PLDataModel, count)
	result := make(chan string, count)

	for w := 0; w < 10; w++ {
		go worker(w, jobs, result)
	}

	for _, m := range models {
		i++
		jobs <- m
	}
	close(jobs)

	step := count / 100
	limit := step
	for ri := 0; ri < count; ri++ {
		<-result
		if ri >= limit{
			toolkit.Printfn("Saving %d of %d (%dpct) in %s", ri, count, ri*100/count,
				time.Since(t0).String())
			limit +=step
		}
	}
	toolkit.Printfn("Done %s", time.Since(t0).String())
}

func worker(wi int, jobs <-chan *gdrj.PLDataModel, r chan<- string) {
	workerConn, _ := modules.GetDboxIConnection("db_godrej")
	defer workerConn.Close()

	for m := range jobs {
		m.ID = m.PrepareID().(string)
		//gdrj.Save(m)
		workerConn.NewQuery().From(m.TableName()).Save().Exec(toolkit.M{}.Set("data", m))
		r <- m.ID
	}
}
