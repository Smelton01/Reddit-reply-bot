deploy:
	docker build -t strimer:multistage -f Dockerfile.multistage .  
	docker tag strimer:multistage us-central1-docker.pkg.dev/tts-server-327215/strimer-bot/strimer
	docker push us-central1-docker.pkg.dev/tts-server-327215/strimer-bot/strimer 