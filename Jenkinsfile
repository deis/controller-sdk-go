def workpath_linux = "/src/github.com/deis/controller-sdk-go"

def gopath_linux = {
	def gopath = pwd() + "/gopath"
	env.GOPATH = gopath
	gopath
}

def workdir_linux = { String gopath ->
	gopath + workpath_linux
}

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
	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath)

	dir(workdir) {
		stage 'Checkout Linux'
			checkout scm
		stage 'Install Linux'
			sh 'make bootstrap'
		stage 'Test Linux'
			sh 'make test'
	}
}

def actual_commit = ''
def git_branch = ''
def go_repo = ''

stage 'Go & Git Info'
node('linux') {

	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath)

	dir(workdir) {
		checkout scm

		// HACK: Recommended approach for getting command output is writing to and then reading a file.
		sh 'mkdir -p tmp'
		sh 'go list . > tmp/GO_LIST'
		// HACK: Jump hurdle to get at actual PR commit rather than merge commit
		sh 'git rev-parse HEAD | git log --pretty=%P -n 1 | cut -c1-40 > tmp/ACTUAL_COMMIT'
		sh 'git describe --all > tmp/GIT_BRANCH'
		go_list = readFile('tmp/GO_LIST').trim()
		actual_commit = readFile('tmp/ACTUAL_COMMIT').trim()
		git_branch = readFile('tmp/GIT_BRANCH').trim()
		// convert 'github.com/deis/controller-sdk-go' to 'github.com/${env.CHANGE_AUTHOR}/controller-sdk-go'
		go_repo = go_list.replace('deis', env.CHANGE_AUTHOR)
	}
}

stage 'Trigger workflow-cli pipeline with this repo and sha'

def parameters = []
if (git_branch != "remotes/origin/master") {
	echo "Passing down SDK_SHA='${actual_commit}' and SDK_GO_REPO='${go_repo}' to the Deis/workflow-cli job..."
	parameters = [
		[$class: 'StringParameterValue', name: 'SDK_SHA', value: actual_commit],
		[$class: 'StringParameterValue', name: 'SDK_GO_REPO', value: go_repo]]
}

build job: 'Deis/workflow-cli/master', parameters: parameters
