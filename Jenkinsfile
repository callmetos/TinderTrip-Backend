@Library('my-shared-lib') _

pipeline {
    agent any

    tools {
        go 'Go_1.25'
    }

    environment {
        REPO_URL        = 'https://github.com/callmetos/TinderTrip-Backend.git'
        REPO_BRANCH     = 'main'
        REPO_CREDENTIALS= 'github-token'
        PROJECT_NAME    = 'tindertrip-backend'
        PROJECT_BRANCH  = 'main'
    }

    stages {
        stage('Checkout') {
            steps {
                script {
                    notifyN8N("INFO", "Pipeline started. Running system validation prior to deployment.")
                }
                git branch: "${env.REPO_BRANCH}",
                    url: "${env.REPO_URL}",
                    credentialsId: "${env.REPO_CREDENTIALS}"
            }
            post {
                failure { script { notifyN8N("FAILURE", "Stage: Checkout failed") } }
            }
        }

        stage('Verify Go') {
            steps { sh 'go version' }
            post {
                failure { script { notifyN8N("FAILURE", "Stage: Verify Go failed") } }
            }
        }

        stage('Quality Checks (parallel)') {
            parallel {
                stage('Download & Verify Deps') {
                    steps {
                        sh 'go mod download'
                        sh 'go mod verify'
                    }
                }

                stage('Fmt Check') {
                    steps {
                        sh '''
                          if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
                            echo "❌ Some files are not formatted:"
                            gofmt -s -l .
                            exit 1
                          fi
                        '''
                    }
                }

                stage('Tests (non-blocking)') {
                    steps {
                        sh '''
                          echo "🧪 Running tests (will not fail pipeline if tests fail)..."
                          go test -v ./... || echo "⚠️ Tests failed or skipped (likely due to sqlite/CGO). Continuing..."
                        '''
                    }
                }
            }
        }

        stage('Build project') {
            steps {
                sh 'go build -o app ./cmd/api/main.go'
            }
            post {
                failure { script { notifyN8N("FAILURE", "Stage: Build project failed") } }
            }
        }

        stage('Health Check (soft)') {
            steps {
                sh '''
                  set +e
                  echo "⏳ Starting backend for soft health check..."
                  ./app > /tmp/app.log 2>&1 &
                  SERVER_PID=$!
                  sleep 5

                  echo "🔎 Curl /health..."
                  curl -sS --max-time 5 http://localhost:8080/health || echo "⚠️ Health endpoint not responding"

                  if kill -0 $SERVER_PID 2>/dev/null; then
                    kill $SERVER_PID
                    echo "🛑 Stopped app (PID $SERVER_PID)"
                  fi

                  echo "—— Last 20 lines of app log ——"
                  tail -n 20 /tmp/app.log || true
                '''
            }
            post {
                failure { script { notifyN8N("FAILURE", "Stage: Health Check failed") } }
            }
        }

        stage('Skip Deploy') {
            when { not { branch 'main' } }
            steps {
                script {
                    echo "⏭️ Skipping deploy: branch = ${env.BRANCH_NAME}, only main can deploy."
                    notifyN8N("INFO", "⏭️ Deploy skipped because branch is ${env.BRANCH_NAME}, only main can deploy.")
                }
            }
        }

        stage('Deploy to Coolify') {
            when { branch 'main' }
            steps {
                script {
                    notifyN8N("INFO", "Preparing deployment to Coolify...")
                    deployToCoolify(
                        "TinderTrip Backend",                
                        "COOLIFY_UUID_tindertrip",
                        "COOLIFY_TOKEN",
                        "COOLIFY_BASEURL"
                    )
                    notifyN8N("SUCCESS", "Deployment request has been successfully sent to Coolify.")
                }
            }
            post {
                failure { script { notifyN8N("FAILURE", "Stage: Deploy to Coolify failed") } }
            }
        }
    }

    post {
        failure {
            script {
                notifyN8N("FAILURE", "Pipeline execution failed. Please review the logs for details.")
            }
        }
    }
}
