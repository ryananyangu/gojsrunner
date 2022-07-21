limit size of file of service
limit function inputs and out puts to specific structures. request build --> response build
passdown statics

cache entire request response templates for fater processing
use grpc for faster inter service communication

--> wrappercore --> invoked script (request, statics,headers) --> string and headers object
--> response --> response payload -(already json) --> structured response


Recieve request
-> Get the configs (constants , headers ,payload) based on the service id
-> Get the client credential if available from the db based on client ID and serviceCode
-> Pass the configs and payload to wrapper script for building including client credentials
-> Get the response from the script
-> send post request to the TP service
-> Get response and convert it to json if not already in json format
-> Pipe to post send request script
-> Get the response back and propagate it to the caller
------------------------------------------------------------------------------------------
