package main

import (
	"eaciit/gdrj/model"
	"eaciit/gdrj/modules"
	"os"
	
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm/v1"
	"github.com/eaciit/toolkit"
    "time"
	//"strings"
)

var conn dbox.IConnection
var count int
var t0 time.Time

var (
    shs = toolkit.M{}
    pcs = toolkit.M{}
    ccs = toolkit.M{}    
    ledgers = toolkit.M{}
    prods = toolkit.M{}
    custs = toolkit.M{}
    vdistskus = toolkit.M{}
)

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

func getCursor(obj orm.IModel)dbox.ICursor{
    c, e := gdrj.Find(obj,nil,nil)
    if e!=nil{
        return nil
    }
    return c
}

func prepMaster(){
    pc:=new(gdrj.ProfitCenter)
    cc:=new(gdrj.CostCenter)
    prod:=new(gdrj.Product)
    ledger:=new(gdrj.LedgerMaster)
    
    cpc := getCursor(pc)
    defer cpc.Close()
    var e error
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
    
    cledger:=getCursor(ledger)
    defer cledger.Close()
    for e=cledger.Fetch(ledger,1,false);e==nil;{
        ledgers.Set(ledger.ID,ledger)
        ledger=new(gdrj.LedgerMaster)
        e=cledger.Fetch(ledger,1,false)
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
}

type progress struct{
    index, count int
}

func worker(id int, jobs <-chan *gdrj.SalesTrx, result chan<- string, p *progress){
    conn, _ := modules.GetDboxIConnection("db_godrej")
    defer conn.Close()
    
    for j := range jobs{
        conn.NewQuery().From("tmpSalesTrx").Save().Exec(toolkit.M{}.Set("data", j))
        p.index++
        log := toolkit.Sprintf("%d of %d Saving: %s in %s", 
            p.index, p.count, j.ID, time.Since(t0).String()) 
        result <- log
    }
}

func main() {
	setinitialconnection()
	defer gdrj.CloseDb()
    
    toolkit.Println("Reading Master")
    //prepMaster()
    
	toolkit.Println("START...")
    seed := 100000
    crx, err := gdrj.Find(new(gdrj.SalesTrx), 
        nil,
        //dbox.Eq("pcvalid",false),
        //nil)
        toolkit.M{}.Set("take",seed))
    if err != nil {
		toolkit.Println("Error Found : ", err.Error())
		os.Exit(1)
	}
    defer crx.Close()
    
    count = seed//crx.Count()
    i := 0
    t0 = time.Now()
    jobs := make(chan *gdrj.SalesTrx, count)
    result := make(chan string)
    
    //--- prepare worker
    p := new(progress)
    p.count = count
    for w:=0;w<5;w++{
        go worker(w, jobs, result, p)
    }
    
    //--- enter jobs
    for {
        i++ 
        st := new(gdrj.SalesTrx)
        e := crx.Fetch(st,1,false)
        if e!=nil {
            break
        }
        
        jobs <- st
        
        toolkit.Printfn("Processing %d of %d %s in %s", 
            i, count, st.SalesHeaderID, 
            time.Since(t0).String())
    }
    close(jobs)
    
    //--- get output
    for i:=0;i<count;i++{
        toolkit.Printfn(<-result)
    }
}