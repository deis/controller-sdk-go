node('windows') {
	def gopath = pwd() + "\\gopath"
	env.GOPATH = gopath
	def workdir = gopath + "\\src\\github.com\\deis\\controller-sdk-go"

	def pscmd = { String cmd ->
		"powershell -NoProfile -ExecutionPolicy Bypass -Command \"${cmd}\""
	}

	dir(workdir) {
		stage 'Checkout Windows'
			checkout scm
		stage 'Install Windows'
			bat pscmd('.\\make bootstrap')
		stage 'Test Windows'
			bat pscmd('.\\make test')
	}
}

node('linux') {
	def gopath = pwd() + "/gopath"
	env.GOPATH = gopath
	def workdir = gopath + "/src/github.com/deis/controller-sdk-go"

	dir(workdir) {
		stage 'Checkout Linux'
			checkout scm
		stage 'Install Linux'
			sh 'make bootstrap'
		stage 'Test Linux'
			sh 'make test'
	}
}
