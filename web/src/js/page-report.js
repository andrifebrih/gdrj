let currentReportMenu = vm.menu().find((d) => d.title === 'Report')
	.submenu.find((d) => d.href == ('/' + document.URL.split('/').slice(3).join('/')))

vm.pageTitle('Report')
vm.breadcrumb([
	{ title: 'Godrej', href: '#' },
	{ title: 'Report', href: '#' },
	{ title: currentReportMenu.title, href: currentReportMenu.href }
])
vm.pageTitle(currentReportMenu.title)

viewModel.report = new Object()
let rpt = viewModel.report

rpt.masterData = {
	// common
	Branch: ko.observableArray(app.repeatAlphabetically('Branch')),
	Brand: ko.observableArray(app.repeatAlphabetically('Brand')),
	SKU: ko.observableArray(app.repeatAlphabetically('SKU')),
	Outlet: ko.observableArray(app.repeatAlphabetically('Outlet')),

	// geo
	Region: ko.observableArray(app.repeatAlphabetically('Region')),
	Area: ko.observableArray(app.repeatAlphabetically('Area')),

	// sales 
	Group: ko.observableArray(app.repeatAlphabetically('Group')),
	KeyAccount: ko.observableArray(app.repeatAlphabetically('Key Account')),
	Channel: ko.observableArray(app.repeatAlphabetically('Channel')),

	// cost
	CC: ko.observableArray(app.repeatAlphabetically('CC')),
	Function: ko.observableArray(app.repeatAlphabetically('Function')),

	// ledger
	Group1: ko.observableArray(app.repeatAlphabetically('Group 1')),
	Group2: ko.observableArray(app.repeatAlphabetically('Group 2')),
	GLAccount: ko.observableArray(app.repeatAlphabetically('GL Account')),
}

rpt.filter = [
	{ _id: 'common', group: 'Common', sub: [
		{ _id: 'Branch', title: 'Branch' },
		{ _id: 'Brand', title: 'Brand' },
		{ _id: 'SKU', title: 'SKU' },
		{ _id: 'Outlet', title: 'Outlet' }
	] },
	{ _id: 'geo', group: 'Geo', sub: [
		{ _id: 'Region', title: 'Region' },
		{ _id: 'Area', title: 'Area' }
	] },
	{ _id: 'sales', group: 'Sales', sub: [
		{ _id: 'Group', title: 'Group' },
		{ _id: 'KeyAccount', title: 'Key Account' },
		{ _id: 'Channel', title: 'Channel' }
	] },
	{ _id: 'cost', group: 'Cost', sub: [
		{ _id: 'CC', title: 'CC' },
		{ _id: 'Function', title: 'Function' }
	] },
	{ _id: 'ledger', group: 'Ledger', sub: [
		{ _id: 'Group1', title: 'Group 1' },
		{ _id: 'Group2', title: 'Group 2' },
		{ _id: 'GLAccount', title: 'GL Account' }
	] },
]