echo "The remediations service kicks off runs against both of our hosts"
echo curl -H "content-type: application/json" -H "Authorization: PSK remediations" -d @examples/remediations.json http://localhost:8000/internal/v2/dispatch
read

cat examples/remediations.json | less
curl -H "content-type: application/json" -H "Authorization: PSK remediations" -d @examples/remediations.json http://localhost:8000/internal/v2/dispatch | jq
read

echo "The configmanager service also kicks off runs against both of our hosts"
echo curl -H "content-type: application/json" -H "Authorization: PSK configmanager" -d @examples/configmanager.json http://localhost:8000/internal/v2/dispatch 
read

cat examples/configmanager.json | less
curl -H "content-type: application/json" -H "Authorization: PSK configmanager" -d @examples/configmanager.json http://localhost:8000/internal/v2/dispatch | jq
read

echo "Now, Sara, the secops engineer for the entire project, can view these runs:"
echo 'SecOps Sara identity: {"identity":{"internal":{"org_id":"aspian"},"account_number":"901578","user":{"username":"sara", "user_id":"sara"},"type":"User"}}'
echo
curl -H "x-rh-identity: eyJpZGVudGl0eSI6eyJpbnRlcm5hbCI6eyJvcmdfaWQiOiJhc3BpYW4ifSwiYWNjb3VudF9udW1iZXIiOiI5MDE1NzgiLCJ1c2VyIjp7InVzZXJuYW1lIjoic2FyYSIsICJ1c2VyX2lkIjoic2FyYSJ9LCJ0eXBlIjoiVXNlciJ9fQo=" http://localhost:8000/api/playbook-dispatcher/v1/runs | jq
read

echo "And David, devops engineer for telemetry, can view this run:"
echo 'DevOps David identity: {"identity":{"internal":{"org_id":"aspian"},"account_number":"901578","user":{"username":"david", "user_id":"david"},"type":"User"}}'
echo
curl -H "x-rh-identity: eyJpZGVudGl0eSI6eyJpbnRlcm5hbCI6eyJvcmdfaWQiOiJhc3BpYW4ifSwiYWNjb3VudF9udW1iZXIiOiI5MDE1NzgiLCJ1c2VyIjp7InVzZXJuYW1lIjoiZGF2aWQiLCAidXNlcl9pZCI6ImRhdmlkIn0sInR5cGUiOiJVc2VyIn19Cg==" http://localhost:8000/api/playbook-dispatcher/v1/runs | jq
read

echo "And Dani, devops engineer for observability, can view this run:"
echo 'DevOps David identity: {"identity":{"internal":{"org_id":"aspian"},"account_number":"901578","user":{"username":"dani", "user_id":"dani"},"type":"User"}}'
echo 
curl -H "x-rh-identity: eyJpZGVudGl0eSI6eyJpbnRlcm5hbCI6eyJvcmdfaWQiOiJhc3BpYW4ifSwiYWNjb3VudF9udW1iZXIiOiI5MDE1NzgiLCJ1c2VyIjp7InVzZXJuYW1lIjoiZGFuaSIsICJ1c2VyX2lkIjoiZGFuaSJ9LCJ0eXBlIjoiVXNlciJ9fQo=" http://localhost:8000/api/playbook-dispatcher/v1/runs | jq
read

