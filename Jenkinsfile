def workpath_linux_root = "/src/github.com/deis"

def make = { String target ->
	try {
		sh "make ${target} fileperms"
	} catch(error) {
		sh 'make fileperms'
		false
	}
}

def gopath_linux = {
	def gopath = pwd() + "/gopath"
	env.GOPATH = gopath
	gopath
}

def workdir_linux = { String gopath, dest ->
	gopath + workpath_linux_root + "/" + dest
}

def sh = { String cmd ->
	wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'XTerm']) {
		sh cmd
	}
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
	def workdir = workdir_linux(gopath, "controller-sdk-go")

	dir(workdir) {
		stage 'Checkout Linux'
			checkout scm
		stage 'Install Linux'
			make 'bootstrap'
		stage 'Test Linux'
			make 'test'
	}
}

def git_commit = ''
def git_branch = ''
def go_repo = ''

stage 'Go & Git Info'
node('linux') {

	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath, "controller-sdk-go")

	dir(workdir) {
		checkout scm

		// HACK: Recommended approach for getting command output is writing to and then reading a file.
		sh 'mkdir -p tmp'
		sh 'git describe --all > tmp/GIT_BRANCH'
		sh 'git rev-parse HEAD > tmp/GIT_COMMIT'
		sh 'go list . > tmp/GO_LIST'
		git_branch = readFile('tmp/GIT_BRANCH').trim()
		git_commit = readFile('tmp/GIT_COMMIT').trim()
		go_repo = readFile('tmp/GO_LIST').trim()

		if (git_branch != "remotes/origin/master") {
			// Determine actual PR commit, if necessary
			sh 'git rev-parse HEAD | git log --pretty=%P -n 1 --date-order > tmp/MERGE_COMMIT_PARENTS'
			sh 'cat tmp/MERGE_COMMIT_PARENTS'
			merge_commit_parents = readFile('tmp/MERGE_COMMIT_PARENTS').trim()
			if (merge_commit_parents.length() > 40) {
				echo 'More than one merge commit parent signifies that the merge commit is not the PR commit'
				echo "Changing git_commit from '${git_commit}' to '${merge_commit_parents.take(40)}'"
				git_commit = merge_commit_parents.take(40)
			} else {
				echo 'Only one merge commit parent signifies that the merge commit is also the PR commit'
				echo "Keeping git_commit as '${git_commit}'"
			}
			// convert 'github.com/deis/controller-sdk-go' to 'github.com/${env.CHANGE_AUTHOR}/controller-sdk-go'
			go_repo = go_repo.replace('deis', env.CHANGE_AUTHOR)
		}
	}
}

stage 'Checkout workflow-cli repo and build/deploy with appropriate updates'
node('linux') {
	def repo = "workflow-cli"
	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath, repo)

	// vars/closures around uploading artifacts to gcs
	def keyfile = "tmp/key.json"

	def getBasePath = { String filepath ->
		def filename = filepath.lastIndexOf(File.separator)
		return filepath.substring(0, filename)
	}

	def upload_artifacts = { String filepath ->
		withCredentials([[$class: 'FileBinding', credentialsId: 'e80fd033-dd76-4d96-be79-6c272726fb82', variable: 'GCSKEY']]) {
			sh "mkdir -p ${getBasePath(filepath)}"
			sh "cat \"\${GCSKEY}\" > ${filepath}"
			make 'upload-gcs'
		}
	}

	dir(workdir) {
		stage "Checkout ${repo}"
		git url: "https://github.com/deis/${repo}.git", branch: "master"

		stage "Build ${repo}"
		if (git_branch != "remotes/origin/master") {
			echo "Skipping build of 386 binaries to shorten CI for Pull Requests"
			env.BUILD_ARCH = "amd64"
		}
		make 'bootstrap'

		stage "Update local glide.yaml with controller-sdk-go repo '${go_repo}' and version '${git_commit}'"

		def pattern = "github\\.com\\/deis\\/controller-sdk-go\\n\\s+version:\\s+[a-f0-9]+"
		def replacement = "${go_repo.replace("/", "\\/")}\\n  version: ${git_commit}"
		sh "perl -i -0pe 's/${pattern}/${replacement}/' glide.yaml"

		def glideYaml = readFile('glide.yaml')
		echo "Updated glide.yaml:\n${glideYaml}"

		make 'glideup'
		sh "VERSION=${git_commit.take(7)} make build-revision"

		stage "Deploy ${repo}"
		upload_artifacts(keyfile)
	}
}

stage 'Trigger e2e tests'
// If build is on master, trigger workflow-test, otherwise, assume build is a PR and trigger workflow-test-pr
waitUntil {
	try {
		def downstreamJob = git_branch == "remotes/origin/master" ? '/workflow-test' : '/workflow-test-pr'
		build job: downstreamJob, parameters: [[$class: 'StringParameterValue', name: 'WORKFLOW_CLI_SHA', value: git_commit]]
		true
	} catch(error) {
			node('linux') {
				if (git_branch != "remotes/origin/master") {
					withCredentials([[$class: 'StringBinding', credentialsId: '8a727911-596f-4057-97c2-b9e23de5268d', variable: 'SLACKEMAIL']]) {
						mail body: """<!DOCTYPE html>
<html>
<head>
<meta content='text/html; charset=UTF-8' http-equiv='Content-Type' />
</head>
<body>
<div>Author: ${env.CHANGE_AUTHOR}<br/>
Branch: ${env.BRANCH_NAME}<br/>
Commit: ${env.CHANGE_TITLE}<br/>
<p><a href="${env.BUILD_URL}console/">Click here</a> to view build logs.</p>
<p><a href="${env.BUILD_URL}input/">Click here</a> to restart e2e.</p>
</div>
</html>
""", from: 'jenkins@ci.deis.io', subject: 'Controller-sdk-go E2E Test Failure', to: env.SLACKEMAIL, mimeType: 'text/html'
					}
					input "Retry the e2e tests?"
				}
			}
		false
	}
}
