#!/bin/bash

# WiFi Billing System - Setup Script

set -e  # Exit on error

echo "WiFi Billing System - Backend Setup"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Check prerequisites
echo "Step 1: Checking prerequisites..."
echo ""

# Check Go
if ! command -v go &> /dev/null; then
    print_error "Go is not installed"
    echo "Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi
print_success "Go installed: $(go version)"

# Check PostgreSQL
if ! command -v psql &> /dev/null; then
    print_warning "PostgreSQL client not found"
    echo "Please ensure PostgreSQL is installed and running"
else
    print_success "PostgreSQL client found"
fi

echo ""
echo "Step 2: Setting up environment..."
echo ""

# Create .env if it doesn't exist
if [ ! -f ".env" ]; then
    print_info "Creating .env file from template..."
    cp .env.example .env
    print_success ".env file created"
    print_warning "Please edit .env with your configuration before running the server"
else
    print_info ".env file already exists"
fi

echo ""
echo "Step 3: Installing Go dependencies..."
echo ""

go mod download
go mod tidy
print_success "Dependencies installed"

echo ""
echo "Step 4: Database setup..."
echo ""

# Read database config from .env
if [ -f ".env" ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

DB_USER=${DB_USER:-postgres}
DB_NAME=${DB_NAME:-wifi_billing}

read -p "Do you want to create/setup the database now? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_info "Creating database $DB_NAME..."
    
    # Try to create database
    if psql -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;" 2>/dev/null; then
        print_success "Database created"
    else
        print_info "Database might already exist (this is fine)"
    fi
    
    print_info "Running migrations..."
    
    # Run migrations using psql
    if psql -U "$DB_USER" -d "$DB_NAME" -f db/migrations/001_init.up.sql > /dev/null 2>&1; then
        print_success "Schema migration completed"
    else
        print_warning "Schema might already be applied"
    fi
    
    if psql -U "$DB_USER" -d "$DB_NAME" -f db/migrations/002_seed.up.sql > /dev/null 2>&1; then
        print_success "Seed data migration completed"
    else
        print_warning "Seed data might already be applied"
    fi
    
    print_success "Database setup completed"
else
    print_info "Skipping database setup"
    print_warning "Remember to run migrations manually:"
    echo "  psql -U $DB_USER -c \"CREATE DATABASE $DB_NAME;\""
    echo "  psql -U $DB_USER -d $DB_NAME -f db/migrations/001_init.up.sql"
    echo "  psql -U $DB_USER -d $DB_NAME -f db/migrations/002_seed.up.sql"
fi

echo ""
echo "Step 5: Building the application..."
echo ""

go build -o bin/wifi-billing-server ./main.go
print_success "Application built successfully"

echo ""
echo "========================================="
print_success "Setup completed successfully!"
echo ""
echo "Next steps:"
echo ""
echo "1. Configure your environment:"
echo "   ${BLUE}nano .env${NC}"
echo ""
echo "2. Required configuration:"
echo "   - Database credentials (DB_USER, DB_PASSWORD, DB_NAME)"
echo "   - M-Pesa API credentials (MPESA_CONSUMER_KEY, MPESA_CONSUMER_SECRET, MPESA_PASSKEY)"
echo "   - JWT secret (JWT_SECRET)"
echo "   - Frontend URL for CORS (FRONTEND_URL)"
echo ""
echo "3. Run the server:"
echo "   ${GREEN}./bin/wifi-billing-server${NC}"
echo "   or"
echo "   ${GREEN}go run main.go${NC}"
echo ""
echo "4. Test the API:"
echo "   ${BLUE}curl http://localhost:8080/health${NC}"
echo ""
echo "5. View migrations status:"
echo "   ${BLUE}psql -U $DB_USER -d $DB_NAME -c \"SELECT * FROM schema_migrations;\"${NC}"
echo ""
echo "Documentation:"
echo "  - README.md - General documentation"
echo "  - MIGRATIONS.md - Database migration guide"
echo "  - .env.example - Configuration reference"
echo ""
print_warning "Don't forget to configure MikroTik router!"
echo "See: docs/MIKROTIK_SETUP.md (if available)"
echo ""
