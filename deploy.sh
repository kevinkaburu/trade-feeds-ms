gcloud container clusters get-credentials ba-frontend-cluster --zone europe-west1-d
docker build -t gcr.io/betika-africa/bet-validator-core:latest . && \
    docker push gcr.io/betika-africa/bet-validator-core:latest  && \
    kubectl patch deployment bet-validator-deployment -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}"