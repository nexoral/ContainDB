# Production Deployment Guide

This guide covers deploying ContainDB in production environments.

## Prerequisites

### System Requirements

- **Operating System**: Ubuntu 20.04+ or Debian 11+
- **Architecture**: x86_64 (amd64)
- **RAM**: Minimum 2GB, Recommended 4GB+
- **Disk Space**: Minimum 10GB free space
- **Network**: Internet connection for pulling Docker images
- **Privileges**: Root/sudo access required

### Software Requirements

- **Docker**: Version 20.10+
- **Docker Compose**: Version 2.0+ (optional, for importing services)
- **Bash**: Version 4.0+
- **Git**: For source installation

## Installation Methods

### Method 1: Automated Installation (Recommended)

```bash
# Download and install ContainDB
curl -fsSL https://raw.githubusercontent.com/nexoral/containdb/main/Scripts/installer.sh | sudo bash -

# Verify installation
sudo containDB --version
```

### Method 2: Debian Package Installation

```bash
# Download the latest .deb package
wget https://github.com/nexoral/containdb/releases/latest/download/containDB_X.X.X_amd64.deb

# Install the package
sudo dpkg -i containDB_X.X.X_amd64.deb

# Fix any dependency issues
sudo apt-get install -f
```

### Method 3: Build from Source

```bash
# Clone the repository
git clone https://github.com/nexoral/containdb.git
cd ContainDB

# Build the binary
chmod +x Scripts/BinBuilder.sh
./Scripts/BinBuilder.sh

# Build the package
chmod +x Scripts/PackageBuilder.sh
./Scripts/PackageBuilder.sh

# Install
sudo dpkg -i Packages/containDB_*.deb
```

## Production Configuration

### Docker Daemon Configuration

Create or edit `/etc/docker/daemon.json`:

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3",
    "labels": "production"
  },
  "storage-driver": "overlay2",
  "userland-proxy": false,
  "live-restore": true,
  "default-ulimits": {
    "nofile": {
      "Name": "nofile",
      "Hard": 64000,
      "Soft": 64000
    }
  }
}
```

Restart Docker:
```bash
sudo systemctl restart docker
```

### Firewall Configuration

```bash
# Enable UFW
sudo ufw enable

# Allow SSH
sudo ufw allow 22/tcp

# Allow only specific IPs to access database ports
# Example: PostgreSQL from application server
sudo ufw allow from 192.168.1.100 to any port 5432 proto tcp

# For development, allow from local network
sudo ufw allow from 192.168.1.0/24 to any port 5432 proto tcp
```

### System Limits

Add to `/etc/security/limits.conf`:

```
*               soft    nofile          65536
*               hard    nofile          65536
*               soft    nproc           4096
*               hard    nproc           4096
```

## Deployment Scenarios

### Scenario 1: Single Server Development Environment

```bash
# Install ContainDB
sudo containDB

# Install required databases
# Example: PostgreSQL + MySQL + Redis
# Select: Install Database → postgresql
# Configure with strong password
# Enable persistence: Yes
# Enable auto-restart: Yes

# Repeat for other databases

# Export configuration for backup
sudo containDB --export
# Saves to docker-compose.yml
```

### Scenario 2: Multi-Environment Setup

```bash
# Development
sudo containDB
# Install databases without persistence
# Use default ports

# Staging
sudo containDB
# Install databases with persistence
# Use strong passwords
# Enable auto-restart

# Production
# Use exported docker-compose.yml from staging
sudo containDB --import staging-compose.yml
```

### Scenario 3: CI/CD Integration

```yaml
# .github/workflows/deploy.yml
name: Deploy Databases

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install ContainDB
        run: |
          curl -fsSL https://raw.githubusercontent.com/nexoral/containdb/main/Scripts/installer.sh | sudo bash -

      - name: Deploy databases
        run: |
          sudo containDB --import docker-compose.yml
```

## High Availability Setup

While ContainDB itself doesn't provide HA clustering, you can achieve high availability:

### Using Docker Swarm

```bash
# Initialize swarm
docker swarm init

# Export ContainDB configuration
sudo containDB --export

# Convert to swarm stack
# Edit docker-compose.yml and add deploy section:
version: '3.8'
services:
  postgres:
    image: postgres
    deploy:
      replicas: 1
      restart_policy:
        condition: on-failure
      placement:
        constraints:
          - node.role == manager
```

### Using Kubernetes

```bash
# Export ContainDB configuration
sudo containDB --export

# Use kompose to convert to Kubernetes manifests
kompose convert -f docker-compose.yml

# Deploy to Kubernetes
kubectl apply -f .
```

## Backup and Recovery

### Backup Strategy

**1. Volume Backup**

```bash
# List ContainDB volumes
docker volume ls | grep "data"

# Backup a volume
docker run --rm \
  -v postgresql-data:/data \
  -v $(pwd):/backup \
  ubuntu tar czf /backup/postgresql-backup-$(date +%Y%m%d).tar.gz /data

# Automated backup script
cat > /usr/local/bin/containdb-backup.sh <<'EOF'
#!/bin/bash
BACKUP_DIR="/var/backups/containdb"
mkdir -p $BACKUP_DIR

for volume in $(docker volume ls -q | grep "\-data$"); do
  echo "Backing up $volume"
  docker run --rm \
    -v $volume:/data \
    -v $BACKUP_DIR:/backup \
    ubuntu tar czf /backup/${volume}-$(date +%Y%m%d-%H%M%S).tar.gz /data
done

# Keep only last 7 days
find $BACKUP_DIR -name "*.tar.gz" -mtime +7 -delete
EOF

chmod +x /usr/local/bin/containdb-backup.sh

# Add to cron
echo "0 2 * * * /usr/local/bin/containdb-backup.sh" | sudo crontab -
```

**2. Export Configuration**

```bash
# Backup ContainDB configuration
sudo containDB --export
cp docker-compose.yml /backup/containdb-config-$(date +%Y%m%d).yml
```

**3. Database-Specific Backups**

```bash
# PostgreSQL
docker exec postgresql-container pg_dumpall -U postgres > backup.sql

# MySQL
docker exec mysql-container mysqldump -u root -p --all-databases > backup.sql

# MongoDB
docker exec mongodb-container mongodump --out=/backup

# Redis
docker exec redis-container redis-cli SAVE
docker cp redis-container:/data/dump.rdb ./redis-backup.rdb
```

### Recovery

**1. Restore Volume**

```bash
# Stop container
docker stop postgresql-container

# Remove old volume
docker volume rm postgresql-data

# Create new volume
docker volume create postgresql-data

# Restore data
docker run --rm \
  -v postgresql-data:/data \
  -v $(pwd):/backup \
  ubuntu tar xzf /backup/postgresql-backup-20260112.tar.gz -C /

# Restart container
docker start postgresql-container
```

**2. Restore from Configuration**

```bash
# Import saved configuration
sudo containDB --import /backup/containdb-config-20260112.yml
```

## Monitoring

### Container Health Monitoring

```bash
# Check container status
docker ps -a --filter "network=ContainDB-Network"

# Monitor resource usage
docker stats

# Check container logs
docker logs -f postgresql-container

# Set up health checks
docker inspect --format='{{json .State.Health}}' postgresql-container
```

### System Monitoring with Prometheus

Create `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'docker'
    static_configs:
      - targets: ['localhost:9323']
```

Enable Docker metrics:

```json
{
  "metrics-addr": "0.0.0.0:9323",
  "experimental": true
}
```

### Log Aggregation

Forward logs to centralized logging:

```bash
# Install Promtail (for Loki)
docker run -d \
  --name promtail \
  --network ContainDB-Network \
  -v /var/log:/var/log \
  -v /var/lib/docker/containers:/var/lib/docker/containers \
  grafana/promtail:latest \
  -config.file=/etc/promtail/config.yml
```

## Performance Optimization

### Docker Performance

```bash
# Clean up unused resources
docker system prune -a --volumes

# Optimize image layers
# Use multi-stage builds for custom images

# Limit container resources
docker update --memory=2g --cpus=2 postgresql-container
```

### Database Performance

**PostgreSQL**
```bash
docker exec postgresql-container bash -c "echo 'shared_buffers = 256MB' >> /var/lib/postgresql/data/postgresql.conf"
docker restart postgresql-container
```

**MySQL**
```bash
docker exec mysql-container bash -c "echo 'innodb_buffer_pool_size = 1G' >> /etc/mysql/my.cnf"
docker restart mysql-container
```

**Redis**
```bash
docker exec redis-container redis-cli CONFIG SET maxmemory 512mb
docker exec redis-container redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

## Security Hardening

### Minimal Privileges

```bash
# Create dedicated user
sudo useradd -r -s /bin/false containdb
sudo usermod -aG docker containdb

# Run with restricted user (note: ContainDB requires sudo)
# Consider using sudo rules to limit commands
```

### Docker Security

```bash
# Scan images for vulnerabilities
docker scan postgres:latest

# Use read-only root filesystem where possible
docker run --read-only --tmpfs /tmp ...

# Drop unnecessary capabilities
docker run --cap-drop=ALL --cap-add=NET_BIND_SERVICE ...
```

### Network Isolation

```bash
# Create separate networks for different environments
docker network create --driver bridge --subnet 172.19.0.0/16 production-db-network
docker network create --driver bridge --subnet 172.20.0.0/16 staging-db-network
```

## Troubleshooting Production Issues

### Container Won't Start

```bash
# Check logs
docker logs postgresql-container

# Check port conflicts
sudo netstat -tulpn | grep 5432

# Check volume permissions
docker run --rm -v postgresql-data:/data ubuntu ls -la /data

# Verify network
docker network inspect ContainDB-Network
```

### Performance Issues

```bash
# Check resource usage
docker stats postgresql-container

# Check I/O wait
iostat -x 5

# Check Docker daemon
sudo systemctl status docker
journalctl -u docker -n 100
```

### Data Corruption

```bash
# Stop container
docker stop postgresql-container

# Run integrity check (PostgreSQL example)
docker run --rm -v postgresql-data:/var/lib/postgresql/data postgres pg_checksums --check -D /var/lib/postgresql/data

# Restore from backup if needed
```

## Updating ContainDB

```bash
# Check current version
sudo containDB --version

# Update via installer
curl -fsSL https://raw.githubusercontent.com/nexoral/containdb/main/Scripts/installer.sh | sudo bash -

# Or download and install new .deb package
wget https://github.com/nexoral/containdb/releases/latest/download/containDB_X.X.X_amd64.deb
sudo dpkg -i containDB_X.X.X_amd64.deb

# Verify update
sudo containDB --version
```

## Migration Guide

### From Native Installation to ContainDB

```bash
# 1. Export native database
pg_dumpall > backup.sql  # PostgreSQL example

# 2. Install ContainDB and create container
sudo containDB
# Select PostgreSQL, configure, install

# 3. Import data
docker cp backup.sql postgresql-container:/tmp/
docker exec -it postgresql-container psql -U postgres -f /tmp/backup.sql

# 4. Verify data
docker exec -it postgresql-container psql -U postgres -c "\l"

# 5. Update application connection strings
# Change: localhost:5432 → localhost:5432 (same if port mapped)
```

### From Docker Compose to ContainDB

```bash
# 1. Export existing setup
docker-compose down
# Save your docker-compose.yml

# 2. Install ContainDB
curl -fsSL https://raw.githubusercontent.com/nexoral/containdb/main/Scripts/installer.sh | sudo bash -

# 3. Import configuration
sudo containDB --import docker-compose.yml

# 4. Verify services
docker ps --filter "network=ContainDB-Network"
```

## Support and Maintenance

### Health Check Script

```bash
#!/bin/bash
# /usr/local/bin/containdb-health.sh

echo "=== ContainDB Health Check ==="
echo "Date: $(date)"
echo

echo "Docker Status:"
systemctl is-active docker

echo -e "\nContainDB Containers:"
docker ps --filter "network=ContainDB-Network" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

echo -e "\nDisk Usage:"
docker system df

echo -e "\nNetwork Status:"
docker network inspect ContainDB-Network --format='{{.Name}}: {{len .Containers}} containers'

echo -e "\nVolume Status:"
docker volume ls --filter "name=-data" --format "table {{.Name}}\t{{.Driver}}\t{{.Mountpoint}}"
```

### Maintenance Schedule

- **Daily**: Automated backups, log rotation
- **Weekly**: Security updates, vulnerability scans
- **Monthly**: Performance review, capacity planning
- **Quarterly**: Disaster recovery drill, update strategy review

---

**Production Checklist:**

- [ ] Docker configured with proper logging
- [ ] Firewall rules configured
- [ ] Automated backups set up
- [ ] Monitoring in place
- [ ] Strong passwords configured
- [ ] SSL/TLS for remote access
- [ ] Regular security updates scheduled
- [ ] Disaster recovery plan tested
- [ ] Documentation updated
- [ ] Team trained on procedures

For additional support, see [TROUBLESHOOTING.md](TROUBLESHOOTING.md) or open an issue on GitHub.
