### Login Request

* Register and enroll new users in Organization - **Org1**:

`curl -s -X POST http://10.53.18.86:4000/users -H "content-type: application/x-www-form-urlencoded" -d 'username=Jim1&orgName=Org1'`

**OUTPUT:**

```
{
  "success": true,
  "secret": "RaxhMgevgJcm",
  "message": "Jim enrolled Successfully",
  "token": "<put JSON Web Token here>"
}
```
Returns Token


### Invoke Hospital request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer1.org1.example.com"],
	"fcn":"hospitalRequest",
	"args":["REQ1","test","syringe","Volunteer1"]
}'
```
Returns Success



### Query Hospital request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer1.org1.example.com"],
	"fcn":"hospitalQuery",
	"args":[]
}'
```
Returns all request made by hospital



### Invoke Volunteer request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org2.example.com","peer1.org2.example.com"],
	"fcn":"volunteerInvoke",
	"args":["REQ1","true","NGO_ORG"]
}'
```
Returns Success



### Query Volunteer request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org2.example.com","peer1.org2.example.com"],
	"fcn":"volunteerQuery",
	"args":[]
}'
```
Returns all request made to Volunteer



### Invoke NGO request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org3.example.com","peer1.org3.example.com"],
	"fcn":"ngoInvoke",
	"args":["REQ1","true"]
}'
```
Returns Success



### Query NGO request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org3.example.com","peer1.org3.example.com"],
	"fcn":"NGOQuery",
	"args":[]
}'
```
Returns all request made to NGO



### Query HOspital Success request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org3.example.com","peer1.org3.example.com"],
	"fcn":"NGOQuery",
	"args":['REQ1','true']
}'
```
Supply chain complete
