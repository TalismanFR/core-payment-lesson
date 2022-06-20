sudo rm -rf vault/data
sudo rm -rf vault/file
sudo rm -rf vault/policies
sudo rm -rf pgdata
sudo rm -rf vault_secrets.txt
docker stop my-id-postgres-1 my-id-vault-1 | xargs docker rm
