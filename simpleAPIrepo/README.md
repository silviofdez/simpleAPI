Folder structure: (go workspace)
	 - bin
	 - pkg
	 - src
	 	- example (mini Angular app mock up using simple API)
	 	- github.com
	 	- simpleAPI (API using go)

In simpleAPI directory you can found the datasets, the source code of the API and a compiled (x86) version.


Go lang was the final choice for the implementation of the API. It supports post and get methods over all the datasets provided. Delete and update are not supported. It is develop using and assuming a Linux systems, in particular it have been developed in Debian 7 32 bits. We can get a working go environment (supposing go and go tools are installed <https://golang.org/doc/install>):

export GOPATH=$HOME/gowksp
export PATH=$PATH:$GOPATH/bin
mkdir $HOME/gowksp/src
cd $HOME/gowksp/src
export SRC=$HOME/gowksp/src

Out API depends on gorilla, so we have to get it:

go get github.com/gorilla/mux

also handlers (to avoid CORS)

go get github.com/gorilla/handlers

Go version >= 1.4. Developed using 1.6.2

You need data en the correct directory tree (dataset/data , dataset/clinical) under the work directory. In the github repository you can find the whole go workspace, as well a binary of the application under the src directory. To run the program just go to the src directory and type:

*go run simpleAPI.go*

Or, if prefer

*go build simpleAPI.go*

./simpleAPI

Examples urls:

\#clinical data dos datasets

curl -X GET http://127.0.0.1:8080/clinicaldata/validation -v

curl -X GET http://127.0.0.1:8080/clinicaldata/validation/UDX201 -v

curl -X POST -H "Content-Type: application/json" -d '{"sample":"UDX555","condition":"C1","age":"99","gender":"Male","race":"Caucasian"}' http://127.0.0.1:8080/clinicaldata/validation/UDX555 -v

\#\#\#

curl -X GET http://127.0.0.1:8080/clinicaldata/training -v

curl -X GET http://127.0.0.1:8080/clinicaldata/training/UDX001 -v

curl -X POST -H "Content-Type: application/json" -d '{"sample":"UDX000","condition":"C1","age":"78","gender":"Male","race":"Other"}' http://127.0.0.1:8080/clinicaldata/training/UDX555 -v

\#\#\#

curl -X GET http://127.0.0.1:8080/data/validation -v

curl -X GET http://127.0.0.1:8080/data/validation/UDX201 -v

curl -X POST -H "Content-Type: application/json" -d '{"sample":"UDX999","m1":"0.3074248313","m2":"0.4074248313","m3":"0.5555248313","m4":"0.3074248313","m5":"0.3074248313","m6":"0.3074248313","m7":"0.3074248313","m8":"0.3074248313","m9":"0.3074248313","m10":"0.3074248313","m11":"0.3074248313","m12":"0.3074248313","m13":"0.3074248313","m14":"0.3074248313","m15":"0.3074248313","m16":"0.3074248313","m17":"0.3074248313","m18":"0.3074248313"}' http://127.0.0.1:8080/data/validation/UDX999 -v

\#\#\#

curl -X GET http://127.0.0.1:8080/data/training -v

curl -X GET http://127.0.0.1:8080/data/training/UDX001 -v

curl -X POST -H "Content-Type: application/json" -d '{"sample":"UDX999","m1":"0.3074248313","m2":"0.4074248313","m3":"0.5555248313","m4":"0.3074248313","m5":"0.3074248313","m6":"0.3074248313","m7":"0.3074248313","m8":"0.3074248313","m9":"0.3074248313","m10":"0.3074248313","m11":"0.3074248313","m12":"0.3074248313","m13":"0.3074248313","m14":"0.3074248313","m15":"0.3074248313","m16":"0.3074248313","m17":"0.3074248313","m18":"0.3074248313"}' http://127.0.0.1:8080/data/training/UDX999 -v

In order to access via the angular app (sorry, very simple mock up, no more time :() you must go to the folder **example** and open index.hmtl using any modern webbrowser from the same machine where the server was.
Data visualization and more complex angular use is on pending tasks.

**NOTE:** This API uses a naïve approach in some aspects: data load, security issues, it is just a proof of concept