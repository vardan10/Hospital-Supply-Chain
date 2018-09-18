# Hospital-NGO-Volunteer Supply Chain
ORG1 - Hospital 1 <br />
ORG2 - NGO 1 <br />
ORG3 - Volunteers <br />

### Register Request

* Register and enroll new users in Organization - **Org1**:

`curl -s -X POST http://10.53.18.86:4000/users -H "content-type: application/x-www-form-urlencoded" -d 'username=Jim&orgName=Org1&secret=enter'`

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


### Login Request

* Enroll new users in Organization - **Org1**:

`curl -s -X POST http://10.53.18.86:4000/loginUsers -H "content-type: application/x-www-form-urlencoded" -d 'username=Jim&orgName=Org1&secret=enter'`

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


##### Do the same for the other 2 orgs.



### Invoke Hospital request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com"],
	"fcn":"hospitalInvoke",
	"args":["REQ1","test","syringe","Volunteer1"]
}'
```
Returns Success



### Query Hospital request

```
curl -s -X GET \
  "http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc?peer=peer0.org1.example.com&fcn=query&args=%5B%27%27%2C%27%27%5D" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
```
Returns all request made by hospital.
args=[''] <- URL ENCODED STRING


### Invoke NGO request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org2.example.com"],
	"fcn":"ngoInvoke",
	"args":["REQ1","volunteer_NAME"]
}'
```
Returns Success



### Query NGO request

```
curl -s -X GET \
  "http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc?peer=peer0.org2.example.com&fcn=query&args=%5B%27%27%2C%27%27%5D" \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json"
```
Returns all request made to NGO
args=[''] <- URL ENCODED STRING


### Invoke Volunteer request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer <put JSON Web Token here>" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org3.example.com"],
	"fcn":"volunteerInvoke",
	"args":["REQ1"]
}'
```
Returns Success



### Query Volunteer request

```
curl -s -X GET \
  "http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc?peer=peer0.org3.example.com&fcn=query&args=%5B%27%27%2C%27%27%5D" \
  -H "authorization: Bearer $ORG3_TOKEN" \
  -H "content-type: application/json"
```
Returns all request made to Volunteer


### Query Hospital Success request

```
curl -s -X POST \
  http://10.53.18.86:4000/channels/mychannel/chaincodes/mycc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com"],
	"fcn":"hospitalSuccess",
	"args":["REQ1"]
}'
```

Supply chain complete
