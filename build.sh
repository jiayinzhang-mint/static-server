NOW=$(date +"%m%d%Y")
VER=$NOW

docker build -t dimlab-docker.pkg.coding.net/insdim/prod/static-server:latest .
docker tag dimlab-docker.pkg.coding.net/insdim/prod/static-server:latest dimlab-docker.pkg.coding.net/insdim/prod/static-server:$VER
docker push dimlab-docker.pkg.coding.net/insdim/prod/static-server:$VER
docker push dimlab-docker.pkg.coding.net/insdim/prod/static-server:latest
