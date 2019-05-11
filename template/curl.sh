curl http://localhost:8080/templates?type=text -v -H "Content-Type:application/json" -d '[
{
	"name":"Deployment",
	"value": {
		"Name":"fist",
		"Image":"sealyun/fist",
		"Replicas":3,
		"Namespace":"sealyun",
		"Command": "['./fist', 'serve']",
		"ImagePolicy":"IfnotPresent",
		"Port":9090}
}
]'

curl http://localhost:8080/templates -v -H "Content-Type:application/json" -d '[
{
	"name":"Deployment",
	"value": {
		"Name":"fist",
		"Image":"sealyun/fist",
		"Replicas":3,
		"Namespace":"sealyun",
		"Command": "['./fist', 'serve']",
		"ImagePolicy":"IfnotPresent",
		"Port":9090}
}
]'