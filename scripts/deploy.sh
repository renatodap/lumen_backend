#!/bin/bash

# LUMEN Backend Deployment Script

set -e

echo "🚀 LUMEN Backend Deployment"
echo "==========================="
echo ""

# Check if platform is specified
PLATFORM=${1:-"docker"}

case $PLATFORM in
  docker)
    echo "🐳 Deploying to Docker..."
    echo ""

    # Build Docker image
    echo "Building Docker image..."
    docker build -t lumen-api:latest .
    echo "✅ Docker image built"
    echo ""

    # Tag for registry (optional)
    if [ -n "$DOCKER_REGISTRY" ]; then
      echo "Tagging image for registry..."
      docker tag lumen-api:latest $DOCKER_REGISTRY/lumen-api:latest
      echo "✅ Image tagged"
      echo ""

      # Push to registry
      echo "Pushing to registry..."
      docker push $DOCKER_REGISTRY/lumen-api:latest
      echo "✅ Image pushed"
    fi

    echo "To run locally:"
    echo "docker run -p 8080:8080 --env-file .env lumen-api:latest"
    ;;

  railway)
    echo "🚂 Deploying to Railway..."
    echo ""

    # Check Railway CLI
    if ! command -v railway &> /dev/null; then
      echo "❌ Railway CLI not installed"
      echo "Install with: npm install -g @railway/cli"
      exit 1
    fi

    # Deploy
    railway up
    echo "✅ Deployed to Railway"
    ;;

  render)
    echo "🎨 Deploying to Render..."
    echo ""

    # Push to git (Render deploys from git)
    git push origin main
    echo "✅ Pushed to git. Render will auto-deploy."
    ;;

  gcp)
    echo "☁️  Deploying to Google Cloud Run..."
    echo ""

    # Check PROJECT_ID
    if [ -z "$GCP_PROJECT_ID" ]; then
      echo "❌ GCP_PROJECT_ID not set"
      exit 1
    fi

    # Build and submit
    gcloud builds submit --tag gcr.io/$GCP_PROJECT_ID/lumen-api

    # Deploy
    gcloud run deploy lumen-api \
      --image gcr.io/$GCP_PROJECT_ID/lumen-api \
      --platform managed \
      --region us-central1 \
      --allow-unauthenticated

    echo "✅ Deployed to Cloud Run"
    ;;

  *)
    echo "❌ Unknown platform: $PLATFORM"
    echo ""
    echo "Usage: ./scripts/deploy.sh [platform]"
    echo ""
    echo "Available platforms:"
    echo "  docker   - Build Docker image"
    echo "  railway  - Deploy to Railway"
    echo "  render   - Deploy to Render"
    echo "  gcp      - Deploy to Google Cloud Run"
    exit 1
    ;;
esac

echo ""
echo "✅ Deployment complete!"
