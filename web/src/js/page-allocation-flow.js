vm.currentMenu('Administration')
vm.currentTitle('Allocation Flow')
vm.breadcrumb([
	{ title: 'Godrej', href: '#' },
	{ title: 'Administration', href: '#' },
	{ title: 'Allocation Flow', href: '/allocationflow' }
])

viewModel.allocationFlow = new Object()
let af = viewModel.allocationFlow

af.templateModuleConfig = {
	Name: '',
	Description: '',
	Params: []
}
af.templateParam = {
	Key: '',
	Value: ''
}
af.moduleConfig = ko.mapping.fromJS(af.templateModuleConfig)
af.isNew = ko.observable(true)
af.dataModules = ko.observableArray([])
af.dataAppliedModules = ko.observableArray([])
af.prepareDrag = () => {
	$('.available-module').sortable({
	    connectWith: '.list-group.module'
	})
	$('.applied-module').sortable({
	    connectWith: '.list-group.module'
	})
}
af.addParam = () => {
	let config = ko.mapping.toJS(af.moduleConfig)
	config.Params.push(app.clone(af.templateParam))
	ko.mapping.fromJS(config, af.moduleConfig)
}
af.removeParam = (index) => {
	return () => {
		let row = af.moduleConfig.Params()[index]
		af.moduleConfig.Params.remove(row)
	}
}
af.getFormPayload = () => {
	let param = new FormData()
	param.append('Name', af.moduleConfig.Name())
	param.append('Description', af.moduleConfig.Description())

	let args = ko.mapping.toJS(af.moduleConfig.Params()).filter((d) => (d.Key != ''))
	param.append('Params', JSON.stringify(args))

	let files = $('[name="file"]')[0].files
	if (files.length > 0) {
		param.append('file', files[0])
	}

	return param
}
af.getModules = () => {
	app.ajaxPost('/allocationflow/getmodules', {}, (res) => {
		if (!app.isFine(res)) {
			return
		}

		af.dataModules(res.data)
		af.prepareDrag()
	})
}
af.getAppliedModules = () => {
	app.ajaxPost('/allocationflow/getappliedmodules', {}, (res) => {
		if (!app.isFine(res)) {
			return
		}

		af.dataAppliedModules(res.data)
		af.prepareDrag()
	})
}
af.doUpload = () => {
	if (!app.isFormValid('.form-upload-file')) {
		return
	}

	let param = af.getFormPayload()
	app.ajaxPost('/allocationflow/uploadnewmodule', param, (res) => {
		if (!app.isFine(res)) {
			return;
		}

		af.getModules()
	})
}
af.showModuleForm = () => {
	app.resetForm($('.form-upload-file'))
	af.addParam()
	$('#modal-module').appendTo($('body'))
	$('#modal-module').modal('show')
}

$(() => {
	af.getModules()
	af.getAppliedModules()
	af.prepareDrag()
})