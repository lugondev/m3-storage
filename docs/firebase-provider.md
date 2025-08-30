# Firebase Storage Provider Configuration

## Overview

Firebase Storage is Google's cloud storage service that integrates seamlessly with Firebase and Google Cloud Platform. It provides secure file uploads and downloads for Firebase apps, with automatic scaling and strong security rules integration.

**When to use Firebase Storage Provider:**
- Applications using Firebase ecosystem (Auth, Firestore, Functions)
- Mobile applications with Firebase SDK integration
- Real-time applications requiring Firebase features
- Applications requiring Firebase Security Rules integration
- Google Cloud Platform integration requirements

**When to consider alternatives:**
- Non-Firebase applications (consider Google Cloud Storage directly)
- Applications hosted on AWS (consider S3)
- Cost-sensitive applications with high bandwidth usage
- Applications requiring advanced S3-compatible features

## Configuration

Add the following configuration to your `config.yaml` file:

```yaml
# Firebase Storage Configuration
firestore:
    projectID: 'your-project-id'                        # Google Cloud Project ID
    credentialsFile: 'path/to/serviceAccountKey.json'   # Service account key file path
    bucketName: 'your-project-id.appspot.com'          # Firebase Storage bucket name
```

## Environment Variables

You can also configure Firebase Storage using environment variables (recommended for production):

- `FIRESTORE_PROJECT_ID`: Google Cloud Project ID
- `FIRESTORE_CREDENTIALS_FILE`: Path to service account key file
- `FIRESTORE_BUCKET_NAME`: Firebase Storage bucket name
- `GOOGLE_APPLICATION_CREDENTIALS`: Google Cloud credentials file path (alternative)

## Configuration Examples

### Production Environment
```yaml
firestore:
    projectID: '${FIREBASE_PROJECT_ID}'                 # Use environment variables
    credentialsFile: '${GOOGLE_APPLICATION_CREDENTIALS}'
    bucketName: 'myapp-prod.appspot.com'
```

### Development Environment
```yaml
firestore:
    projectID: 'myapp-dev-12345'
    credentialsFile: './config/firebase-dev-key.json'
    bucketName: 'myapp-dev-12345.appspot.com'
```

### Multiple Environment Setup
```yaml
# config.production.yaml
firestore:
    projectID: 'myapp-prod'
    credentialsFile: '/secrets/firebase-prod-key.json'
    bucketName: 'myapp-prod.appspot.com'

# config.staging.yaml
firestore:
    projectID: 'myapp-staging'
    credentialsFile: '/secrets/firebase-staging-key.json'
    bucketName: 'myapp-staging.appspot.com'
```

### Docker Container Setup
```yaml
firestore:
    projectID: '${FIREBASE_PROJECT_ID}'
    credentialsFile: '/app/secrets/firebase-key.json'   # Mount secret as volume
    bucketName: '${FIREBASE_PROJECT_ID}.appspot.com'
```

### Custom Bucket Configuration
```yaml
firestore:
    projectID: 'myapp-12345'
    credentialsFile: './config/serviceAccountKey.json'
    bucketName: 'custom-bucket-name'                    # Custom bucket instead of default
```

## Service Account Setup

### Creating Service Account

1. **Go to Google Cloud Console**: https://console.cloud.google.com/
2. **Navigate to IAM & Admin > Service Accounts**
3. **Create Service Account** with appropriate permissions
4. **Generate and download JSON key file**

### Required IAM Roles

Assign these roles to your service account:

```json
[
  "roles/storage.admin",           // Full storage access
  "roles/storage.objectAdmin",     // Object-level access (alternative)
  "roles/firebase.admin"           // Firebase admin access
]
```

### Service Account Key Example
```json
{
  "type": "service_account",
  "project_id": "your-project-id",
  "private_key_id": "key-id",
  "private_key": "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
  "client_email": "your-service-account@your-project-id.iam.gserviceaccount.com",
  "client_id": "123456789",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/your-service-account%40your-project-id.iam.gserviceaccount.com"
}
```

## Features

The Firebase Storage provider supports all standard storage operations:

- **Upload**: Upload files to Firebase Storage buckets
- **Download**: Download files from Firebase Storage
- **Delete**: Remove files from Firebase Storage buckets
- **GetURL**: Get public download URLs for files
- **GetSignedURL**: Generate signed URLs for secure, time-limited access
- **GetObject**: Retrieve file metadata and properties
- **CheckHealth**: Verify connection to Firebase Storage and bucket access

## Firebase Storage Security Rules

Configure Security Rules for client-side access:

### Basic Rules Example
```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    // Allow read/write for authenticated users
    match /{allPaths=**} {
      allow read, write: if request.auth != null;
    }
  }
}
```

### User-specific Rules
```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    // Users can only access their own files
    match /users/{userId}/{allPaths=**} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
    
    // Public read access for certain paths
    match /public/{allPaths=**} {
      allow read;
      allow write: if request.auth != null;
    }
  }
}
```

### Advanced Rules with Validation
```javascript
rules_version = '2';
service firebase.storage {
  match /b/{bucket}/o {
    match /images/{userId}/{imageId} {
      allow read: if true;  // Public read
      allow write: if request.auth != null 
                   && request.auth.uid == userId
                   && resource.size < 5 * 1024 * 1024  // Max 5MB
                   && request.resource.contentType.matches('image/.*');
    }
  }
}
```

## Storage Structure

Organize files in Firebase Storage:

```
gs://your-bucket-name/
├── users/
│   ├── {userId}/
│   │   ├── profile/
│   │   │   └── avatar.jpg
│   │   ├── documents/
│   │   │   └── document.pdf
│   │   └── media/
│   │       ├── image1.png
│   │       └── video1.mp4
├── public/
│   ├── assets/
│   └── shared/
└── temp/
    └── uploads/
```

## Requirements

- Google Cloud Project with Firebase enabled
- Firebase Storage bucket (automatically created with Firebase project)
- Service Account with appropriate permissions
- Service Account key file (JSON format)
- Network connectivity to Firebase/Google Cloud endpoints

## Authentication Methods

### Service Account Key (Current Implementation)
```yaml
firestore:
    projectID: 'your-project-id'
    credentialsFile: './path/to/serviceAccountKey.json'
    bucketName: 'your-project-id.appspot.com'
```

### Application Default Credentials (Google Cloud environments)
```bash
# Set environment variable
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/serviceAccountKey.json"
```

```yaml
firestore:
    projectID: 'your-project-id'
    bucketName: 'your-project-id.appspot.com'
    # credentialsFile can be omitted when using default credentials
```

## Security Considerations

1. **Service Account Security**:
   - Store service account keys securely
   - Use environment variables for credentials
   - Rotate service account keys regularly
   - Follow principle of least privilege

2. **Storage Security Rules**:
   - Implement proper authentication checks
   - Validate file types and sizes
   - Restrict access based on user roles
   - Test rules thoroughly before deployment

3. **Network Security**:
   - Use HTTPS for all requests
   - Configure proper CORS settings
   - Implement rate limiting where needed

4. **Data Protection**:
   - Enable audit logging
   - Use versioning for important files
   - Implement backup strategies

## Performance Optimization

1. **File Organization**: Use logical folder structures for better organization
2. **Caching**: Implement client-side caching for frequently accessed files
3. **CDN**: Firebase automatically provides global CDN
4. **Parallel Uploads**: Support for concurrent uploads
5. **Resumable Uploads**: Automatic resumable upload for large files

## Cost Optimization

1. **Storage Classes**: Firebase Storage uses Google Cloud Storage classes
2. **Lifecycle Policies**: Set up automatic deletion of temporary files
3. **Compression**: Compress files before upload when possible
4. **Monitoring**: Use Firebase Console to monitor usage and costs

### Lifecycle Management
```yaml
# Example lifecycle configuration
lifecycle:
  rule:
    - action:
        type: Delete
      condition:
        age: 30  # Delete files older than 30 days
        matchesPrefix: ['temp/']
```

## Troubleshooting

### Common Issues

1. **Authentication Failed**:
   - Verify service account key file path and format
   - Check if service account has proper permissions
   - Ensure project ID matches the service account project
   - Verify Firebase is enabled for the project

2. **Bucket Not Found**:
   - Verify bucket name format (usually project-id.appspot.com)
   - Ensure Firebase Storage is enabled for the project
   - Check if custom bucket exists

3. **Permission Denied**:
   - Verify IAM roles assigned to service account
   - Check Firebase Security Rules (for client access)
   - Ensure proper API enablement

4. **Network Issues**:
   - Check internet connectivity
   - Verify firewall settings for Google Cloud endpoints
   - Check proxy configuration if applicable

### Health Check

Verify Firebase Storage connection using the health check endpoint:

```bash
curl -X GET "http://localhost:8083/api/v1/storage/health?provider_type=firebase"
```

### Firebase CLI Testing

Test your Firebase Storage configuration using Firebase CLI:

```bash
# Install Firebase CLI
npm install -g firebase-tools

# Login to Firebase
firebase login

# Set project
firebase use your-project-id

# Deploy security rules
firebase deploy --only storage

# Test upload
firebase storage:upload test.txt gs://your-bucket/test.txt

# Test download
firebase storage:download gs://your-bucket/test.txt ./downloaded-test.txt
```

### Google Cloud CLI Testing

Test using gcloud CLI:

```bash
# Authenticate with service account
gcloud auth activate-service-account --key-file=path/to/serviceAccountKey.json

# Set project
gcloud config set project your-project-id

# List buckets
gsutil ls

# Upload file
gsutil cp test.txt gs://your-bucket/test.txt

# Download file
gsutil cp gs://your-bucket/test.txt ./downloaded-test.txt
```

## Monitoring and Logging

### Firebase Console
Monitor usage through Firebase Console:
- Storage usage and bandwidth
- Request statistics
- Security rules debugging
- Performance monitoring

### Google Cloud Monitoring
Set up monitoring for:
- Storage operations
- Request latency
- Error rates
- Bandwidth usage

### Logging
Enable audit logging for security and compliance:
```yaml
auditConfigs:
- service: storage.googleapis.com
  auditLogConfigs:
  - logType: ADMIN_READ
  - logType: DATA_READ
  - logType: DATA_WRITE
```

## Firebase Emulator (Development)

Use Firebase Emulator for local development:

```bash
# Install Firebase CLI
npm install -g firebase-tools

# Initialize Firebase project
firebase init

# Start emulators
firebase emulators:start --only storage

# Storage emulator runs on http://localhost:9199
```

```yaml
# Development configuration with emulator
firestore:
    projectID: 'demo-project'
    bucketName: 'demo-project.appspot.com'
    emulatorHost: 'localhost:9199'  # Use emulator endpoint
```

## Provider Type

When using the storage factory or API endpoints, use provider type: `"firebase"`

## Migration

### From Other Providers to Firebase Storage

1. **Firebase Project Setup**: Create or configure Firebase project
2. **Service Account Setup**: Create service account with proper permissions
3. **Bucket Configuration**: Configure Firebase Storage bucket
4. **Data Transfer**: Use gsutil, custom scripts, or transfer services
5. **Security Rules**: Configure appropriate security rules
6. **Application Update**: Update configuration to use Firebase provider
7. **Testing**: Test functionality thoroughly
8. **Cutover**: Switch to Firebase Storage

### Migration with gsutil
```bash
# Copy from AWS S3 to Firebase Storage
gsutil -m cp -r s3://source-bucket/* gs://firebase-bucket/

# Copy from local storage to Firebase Storage
gsutil -m cp -r ./uploads/* gs://firebase-bucket/
```

## Integration with Firebase Services

### Firebase Authentication
```javascript
// Client-side integration with Firebase Auth
const user = firebase.auth().currentUser;
if (user) {
  // User is signed in, can access storage
  const storageRef = firebase.storage().ref();
  const userRef = storageRef.child(`users/${user.uid}/`);
}
```

### Cloud Functions Integration
```javascript
// Cloud Functions trigger on storage upload
const functions = require('firebase-functions');
const admin = require('firebase-admin');

exports.processUpload = functions.storage.object().onFinalize(async (object) => {
  console.log('File uploaded:', object.name);
  // Process the uploaded file
});
```

### Firestore Integration
```javascript
// Store file metadata in Firestore
const db = firebase.firestore();
const metadata = {
  fileName: 'image.jpg',
  uploadedBy: user.uid,
  timestamp: firebase.firestore.FieldValue.serverTimestamp(),
  storageUrl: downloadURL
};

await db.collection('media').add(metadata);
```