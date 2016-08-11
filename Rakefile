# Copyright 2016 Claudemiro Alves Feitosa Neto. All rights reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.

require 'rake/clean'

VERSION = 'v1.1.0'
GITHASH = `git rev-parse --short HEAD`
DATE = Time.now.strftime '%Y%m%d%H%M%S'

CLOBBER.include 'build'

task :default => [:'run-debug']

desc 'Build a debug version'
task :debug do
	sh "GO15VENDOREXPERIMENT=1 go install -ldflags '-w -X main.version=DEBUG -X main.buildstamp=DEBUG -X main.githash=DEBUG' github.com/dimiro1/ipe"
end

desc 'Build and run debug version'
task :'run-debug' => :debug do
	sh '$GOPATH/bin/ipe --config $GOPATH/src/github.com/dimiro1/ipe/config.json -logtostderr=true -v=2'
end

desc 'Run test suite'
task :test do
	sh 'GO15VENDOREXPERIMENT=1 go test . `glide nv`'
end

desc 'Download the dependencies'
task :'deps' do
	sh 'glide install -v -s'
end

desc 'Generate distributions'
task :distribute => [:linux, :darwin]

desc 'Generate a linux distribution'
task :linux do
	Rake::Task['build'].invoke 'linux'
end

desc 'Generate a darwin distribution'
task :darwin do
	Rake::Task['build'].invoke 'darwin'
end

task :build, [:os] do |t, args|
	t.reenable
	os = args[:os]

	sh "mkdir -p build/#{os}"
	sh "GO15VENDOREXPERIMENT=1 GOOS=#{os} GOARCH=amd64 go build -ldflags '-X main.version=#{VERSION} -X main.buildstamp=#{DATE} -X main.githash=#{GITHASH}' -o build/#{os}/ipe github.com/dimiro1/ipe"
	sh "cp ipe/config-example.json build/#{os}/config.json"
	sh "cp LICENSE build/#{os}/"
	sh "cp README.md build/#{os}/"
	sh "tar -C build/#{os} -czf build/ipe_#{VERSION}_#{os}_amd64.tar.gz ."
end