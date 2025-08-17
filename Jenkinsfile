pipeline {
    agent any
    
    environment {
        DOCKER_HOST = 'ssh://debian@51.178.17.78'
        DEPLOY_HOST = '51.178.17.78'
        DEPLOY_USER = 'debian'
        APP_NAME = 'lol-fullstack-app'
        DOCKER_COMPOSE_FILE = 'docker-compose.complete.yml'
    }
    
    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        timeout(time: 30, unit: 'MINUTES')
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out source code...'
                checkout scm
            }
        }
        
        stage('Build & Test') {
            parallel {
                stage('Backend Build') {
                    steps {
                        echo 'Building Go backend...'
                        sh '''
                            go version
                            go mod tidy
                            go build -o main .
                            go test ./...
                        '''
                    }
                }
                
                stage('Frontend Build') {
                    steps {
                        echo 'Building React frontend...'
                        dir('web') {
                            sh '''
                                if [ -f package.json ]; then
                                    npm install
                                    npm run build
                                else
                                    echo "Frontend already built or not required"
                                fi
                            '''
                        }
                    }
                }
            }
        }
        
        stage('Docker Build') {
            steps {
                echo 'Building Docker image...'
                sh '''
                    docker build -f Dockerfile.simple-fullstack -t lol-match-exporter-fullstack:${BUILD_NUMBER} .
                    docker tag lol-match-exporter-fullstack:${BUILD_NUMBER} lol-match-exporter-fullstack:latest
                '''
            }
        }
        
        stage('Security Scan') {
            steps {
                echo 'Running security scans...'
                sh '''
                    # Scan des vulnérabilités Docker
                    if command -v trivy &> /dev/null; then
                        trivy image lol-match-exporter-fullstack:latest
                    else
                        echo "Trivy not installed, skipping security scan"
                    fi
                    
                    # Scan des dépendances Go
                    if command -v govulncheck &> /dev/null; then
                        govulncheck ./...
                    else
                        echo "govulncheck not installed, skipping Go vulnerability check"
                    fi
                '''
            }
        }
        
        stage('Deploy to Staging') {
            when {
                not { branch 'main' }
            }
            steps {
                echo 'Deploying to staging environment...'
                deployToEnvironment('staging')
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            steps {
                echo 'Deploying to production environment...'
                script {
                    // Demande de confirmation pour la production
                    input message: 'Deploy to production?', ok: 'Deploy',
                          submitterParameter: 'SUBMITTER'
                }
                deployToEnvironment('production')
            }
        }
        
        stage('Health Check') {
            steps {
                echo 'Performing health checks...'
                sh '''
                    # Attendre que l'application démarre
                    sleep 30
                    
                    # Test de santé de l'API
                    curl -f https://herald.lol/api/health || exit 1
                    
                    # Test des endpoints critiques
                    curl -f https://herald.lol/api/auth/session || exit 1
                    curl -f https://herald.lol/api/auth/regions || exit 1
                    
                    # Test de validation auth
                    curl -X POST -H "Content-Type: application/json" \
                         -d '{"username":"test","tagline":"test","region":"euw1"}' \
                         https://herald.lol/api/auth/validate || exit 1
                    
                    echo "All health checks passed!"
                '''
            }
        }
        
        stage('Performance Test') {
            steps {
                echo 'Running performance tests...'
                sh '''
                    # Test de charge basique avec Apache Bench
                    if command -v ab &> /dev/null; then
                        ab -n 100 -c 10 https://herald.lol/api/health
                    else
                        echo "Apache Bench not available, skipping performance test"
                    fi
                '''
            }
        }
    }
    
    post {
        always {
            echo 'Cleaning up...'
            sh '''
                # Nettoyage des images Docker anciennes
                docker image prune -f --filter "until=24h"
            '''
        }
        
        success {
            echo 'Deployment successful!'
            // Notifications Slack/Teams/Email
            script {
                if (env.SLACK_WEBHOOK) {
                    sh """
                        curl -X POST -H 'Content-type: application/json' \
                             --data '{"text":"✅ LoL Match Exporter deployed successfully to production!\\nBuild: ${BUILD_NUMBER}\\nCommit: ${GIT_COMMIT}"}' \
                             ${SLACK_WEBHOOK}
                    """
                }
            }
        }
        
        failure {
            echo 'Deployment failed!'
            // Rollback automatique en cas d'échec
            script {
                if (env.BRANCH_NAME == 'main') {
                    echo 'Rolling back production deployment...'
                    rollbackProduction()
                }
            }
            
            // Notifications d'échec
            script {
                if (env.SLACK_WEBHOOK) {
                    sh """
                        curl -X POST -H 'Content-type: application/json' \
                             --data '{"text":"❌ LoL Match Exporter deployment FAILED!\\nBuild: ${BUILD_NUMBER}\\nBranch: ${BRANCH_NAME}\\nCheck Jenkins for details."}' \
                             ${SLACK_WEBHOOK}
                    """
                }
            }
        }
        
        unstable {
            echo 'Build unstable - some tests may have failed'
        }
    }
}

def deployToEnvironment(environment) {
    sh """
        echo "Deploying to ${environment} environment..."
        
        # Sauvegarder l'image Docker sur le serveur
        docker save lol-match-exporter-fullstack:latest | \
            ssh ${DEPLOY_USER}@${DEPLOY_HOST} 'docker load'
        
        # Copier les fichiers de configuration
        scp ${DOCKER_COMPOSE_FILE} ${DEPLOY_USER}@${DEPLOY_HOST}:~/docker-compose.yml
        scp nginx/nginx-fullstack.conf ${DEPLOY_USER}@${DEPLOY_HOST}:~/nginx.conf
        scp scripts/start-fullstack.sh ${DEPLOY_USER}@${DEPLOY_HOST}:~/start.sh
        
        # Déployer sur le serveur distant
        ssh ${DEPLOY_USER}@${DEPLOY_HOST} '''
            # Arrêter l'ancienne version
            docker-compose down --remove-orphans || true
            
            # Nettoyer les anciens conteneurs
            docker container prune -f
            
            # Démarrer la nouvelle version
            chmod +x start.sh
            ./start.sh
            
            # Vérifier que le conteneur démarre
            sleep 10
            docker ps | grep lol-fullstack-app || exit 1
            
            echo "Deployment to ${environment} completed successfully"
        '''
    """
}

def rollbackProduction() {
    sh """
        echo "Rolling back production to previous version..."
        
        ssh ${DEPLOY_USER}@${DEPLOY_HOST} '''
            # Récupérer la dernière version stable
            LAST_STABLE=\$(docker images --format "table {{.Repository}}:{{.Tag}}" | \
                          grep lol-match-exporter-fullstack | \
                          grep -v latest | \
                          head -1)
            
            if [ -n "\$LAST_STABLE" ]; then
                echo "Rolling back to \$LAST_STABLE"
                
                # Arrêter la version actuelle
                docker-compose down
                
                # Mettre à jour le tag latest vers la version stable
                docker tag \$LAST_STABLE lol-match-exporter-fullstack:latest
                
                # Redémarrer
                ./start.sh
                
                echo "Rollback completed to \$LAST_STABLE"
            else
                echo "No previous version found for rollback"
                exit 1
            fi
        '''
    """
}
