

const API_BASE_URL = window.location.origin;

// DOM Elements
const statusIndicator = document.getElementById('status-indicator');
const statusText = document.getElementById('status-text');
const bucketList = document.getElementById('bucket-list');
const newBucketInput = document.getElementById('new-bucket');
const createBucketBtn = document.getElementById('create-bucket');
const currentBucketHeader = document.getElementById('current-bucket');
const filesList = document.getElementById('files-list');
const fileInput = document.getElementById('file-input');
const customKeyInput = document.getElementById('custom-key');
const uploadBtn = document.getElementById('upload-btn');
const uploadForm = document.getElementById('upload-form');
const dropArea = document.getElementById('drop-area');
const refreshFilesBtn = document.getElementById('refresh-files');
const toastContainer = document.getElementById('toast-container');

// State
let currentBucket = '';
let selectedFile = null;

// Initialize the application
document.addEventListener('DOMContentLoaded', () => {
    checkServerStatus();
    loadBuckets();
    setupEventListeners();
});

// Setup event listeners
function setupEventListeners() {
    // File selection
    fileInput.addEventListener('change', handleFileSelect);
    
    // Form submission
    uploadForm.addEventListener('submit', handleFileUpload);
    
    // Create bucket button
    createBucketBtn.addEventListener('click', handleCreateBucket);
    
    // Create bucket on Enter key in input
    newBucketInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleCreateBucket();
        }
    });
    
    // Drag and drop
    dropArea.addEventListener('dragover', (e) => {
        e.preventDefault();
        dropArea.classList.add('drag-over');
    });
    
    dropArea.addEventListener('dragleave', () => {
        dropArea.classList.remove('drag-over');
    });
    
    dropArea.addEventListener('drop', (e) => {
        e.preventDefault();
        dropArea.classList.remove('drag-over');
        
        if (e.dataTransfer.files.length) {
            fileInput.files = e.dataTransfer.files;
            handleFileSelect({ target: fileInput });
        }
    });
    
    // Refresh files
    refreshFilesBtn.addEventListener('click', () => {
        if (currentBucket) {
            loadFiles(currentBucket);
        }
    });
    
    // Preview modal events
    closePreviewBtn.addEventListener('click', closePreview);
    previewModal.addEventListener('click', (e) => {
        if (e.target === previewModal) {
            closePreview();
        }
    });
    
    // Close preview on Escape key
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape' && previewModal.classList.contains('active')) {
            closePreview();
        }
    });
}

// Check server status
async function checkServerStatus() {
    try {
        const response = await fetch(`${API_BASE_URL}/health`);
        if (response.ok) {
            statusIndicator.classList.add('online');
            statusText.textContent = 'Server Online';
        } else {
            throw new Error('Server returned an error');
        }
    } catch (error) {
        statusIndicator.classList.add('offline');
        statusText.textContent = 'Server Offline';
        showToast('Cannot connect to the S3 Clone server', 'error');
    }
}

// Load all buckets
async function loadBuckets() {
    try {
        const response = await fetch(`${API_BASE_URL}/buckets`);
        if (!response.ok) throw new Error('Failed to load buckets');
        
        const data = await response.json();
        renderBuckets(data.buckets);
    } catch (error) {
        showToast('Failed to load buckets', 'error');
        bucketList.innerHTML = '<li class="loading">Failed to load buckets</li>';
    }
}

// Render buckets in the sidebar
function renderBuckets(buckets) {
    if (!buckets || buckets.length === 0) {
        bucketList.innerHTML = '<li class="loading">No buckets found</li>';
        return;
    }
    
    bucketList.innerHTML = '';
    buckets.forEach(bucket => {
        const li = document.createElement('li');
        
        // Create bucket name span (clickable)
        const nameSpan = document.createElement('span');
        nameSpan.className = 'bucket-name';
        nameSpan.textContent = bucket;
        nameSpan.addEventListener('click', () => selectBucket(bucket));
        
        // Create delete button
        const deleteBtn = document.createElement('span');
        deleteBtn.className = 'bucket-delete';
        deleteBtn.innerHTML = '<i class="fas fa-trash"></i>';
        deleteBtn.title = 'Delete bucket';
        deleteBtn.addEventListener('click', (e) => {
            e.stopPropagation(); // Prevent bucket selection when clicking delete
            deleteBucket(bucket);
        });
        
        // Add elements to li
        li.appendChild(nameSpan);
        li.appendChild(deleteBtn);
        
        if (bucket === currentBucket) {
            li.classList.add('active');
        }
        
        bucketList.appendChild(li);
    });
}

// Select a bucket
function selectBucket(bucket) {
    currentBucket = bucket;
    currentBucketHeader.textContent = `Bucket: ${bucket}`;
    
    // Update active bucket in the sidebar
    const bucketItems = bucketList.querySelectorAll('li');
    bucketItems.forEach(item => {
        if (item.textContent === bucket) {
            item.classList.add('active');
        } else {
            item.classList.remove('active');
        }
    });
    
    // Enable upload button if a file is selected
    if (selectedFile) {
        uploadBtn.disabled = false;
    }
    
    // Load files in the bucket
    loadFiles(bucket);
}

// Load files in a bucket
async function loadFiles(bucket) {
    try {
        refreshFilesBtn.querySelector('i').classList.add('loading');
        const response = await fetch(`${API_BASE_URL}/list/${bucket}`);
        
        if (!response.ok) throw new Error('Failed to load files');
        
        const data = await response.json();
        renderFiles(data.files);
    } catch (error) {
        showToast('Failed to load files', 'error');
        filesList.innerHTML = '<div class="no-files">Failed to load files</div>';
    } finally {
        refreshFilesBtn.querySelector('i').classList.remove('loading');
    }
}

// Render files in the selected bucket
function renderFiles(files) {
    if (!files || files.length === 0) {
        filesList.innerHTML = '<div class="no-files">No files in this bucket</div>';
        return;
    }
    
    filesList.innerHTML = '';
    files.forEach(file => {
        const fileItem = document.createElement('div');
        fileItem.className = 'file-item';
        
        // Get file extension for icon
        const extension = file.split('.').pop().toLowerCase();
        let iconClass = 'fa-file';
        
        // Choose appropriate icon based on file type
        if (['jpg', 'jpeg', 'png', 'gif', 'svg', 'webp'].includes(extension)) {
            iconClass = 'fa-file-image';
        } else if (['pdf'].includes(extension)) {
            iconClass = 'fa-file-pdf';
        } else if (['doc', 'docx'].includes(extension)) {
            iconClass = 'fa-file-word';
        } else if (['xls', 'xlsx', 'csv'].includes(extension)) {
            iconClass = 'fa-file-excel';
        } else if (['zip', 'rar', 'tar', 'gz'].includes(extension)) {
            iconClass = 'fa-file-archive';
        }
        
        fileItem.innerHTML = `
            <div class="file-info">
                <div class="file-icon">
                    <i class="fas ${iconClass}"></i>
                </div>
                <div class="file-name">${file}</div>
            </div>
            <div class="file-actions">
                <div class="file-action download" title="Download">
                    <i class="fas fa-download"></i>
                </div>
                <div class="file-action delete" title="Delete">
                    <i class="fas fa-trash"></i>
                </div>
            </div>
        `;
        
        // Add event listeners
        const downloadBtn = fileItem.querySelector('.download');
        downloadBtn.addEventListener('click', () => downloadFile(file));
        
        const deleteBtn = fileItem.querySelector('.delete');
        deleteBtn.addEventListener('click', () => deleteFile(file));
        
        filesList.appendChild(fileItem);
    });
}

// Handle file selection
function handleFileSelect(event) {
    selectedFile = event.target.files[0];
    
    if (selectedFile) {
        const fileMsg = dropArea.querySelector('.file-msg');
        fileMsg.innerHTML = `
            <i class="fas fa-file"></i>
            <span class="file-selected">${selectedFile.name} (${formatFileSize(selectedFile.size)})</span>
        `;
        
        // Enable upload button if bucket is selected
        if (currentBucket) {
            uploadBtn.disabled = false;
        }
    }
}

// Format file size
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Handle file upload
async function handleFileUpload(event) {
    event.preventDefault();
    
    if (!selectedFile || !currentBucket) {
        showToast('Please select a file and a bucket', 'error');
        return;
    }
    
    const formData = new FormData();
    formData.append('file', selectedFile);
    
    let url = `${API_BASE_URL}/upload/${currentBucket}`;
    const customKey = customKeyInput.value.trim();
    
    // Add custom key as query parameter if provided
    if (customKey) {
        url += `?key=${encodeURIComponent(customKey)}`;
    }
    
    try {
        uploadBtn.disabled = true;
        uploadBtn.textContent = 'Uploading...';
        
        const response = await fetch(url, {
            method: 'PUT',
            body: formData
        });
        
        if (!response.ok) {
            throw new Error('Failed to upload file');
        }
        
        const data = await response.json();
        showToast(`File uploaded successfully as ${data.key}`, 'success');
        
        // Reset form
        uploadForm.reset();
        customKeyInput.value = '';
        selectedFile = null;
        uploadBtn.disabled = true;
        uploadBtn.textContent = 'Upload';
        
        // Reset file message
        const fileMsg = dropArea.querySelector('.file-msg');
        fileMsg.innerHTML = `
            <i class="fas fa-cloud-upload-alt"></i>
            <span>Drag & drop or click to upload files</span>
        `;
        
        // Reload files
        loadFiles(currentBucket);
    } catch (error) {
        showToast('Failed to upload file', 'error');
        uploadBtn.disabled = false;
    } finally {
        uploadBtn.textContent = 'Upload';
    }
}

// Download a file
function downloadFile(filename) {
    if (!currentBucket) return;
    
    const downloadUrl = `${API_BASE_URL}/download/${currentBucket}/${filename}`;
    window.open(downloadUrl, '_blank');
}

// Delete a file
async function deleteFile(filename) {
    if (!currentBucket) return;
    
    if (!confirm(`Are you sure you want to delete ${filename}?`)) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/delete/${currentBucket}/${filename}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) {
            throw new Error('Failed to delete file');
        }
        
        showToast(`File ${filename} deleted successfully`, 'success');
        loadFiles(currentBucket);
    } catch (error) {
        showToast('Failed to delete file', 'error');
    }
}

// Handle bucket creation
async function handleCreateBucket() {
    const bucketName = newBucketInput.value.trim();
    
    if (!bucketName) {
        showToast('Please enter a bucket name', 'error');
        return;
    }
    
    try {
        createBucketBtn.disabled = true;
        
        const response = await fetch(`${API_BASE_URL}/buckets`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name: bucketName })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showToast(`Bucket '${bucketName}' created successfully`, 'success');
            newBucketInput.value = '';
            loadBuckets();
        } else {
            showToast(data.error || 'Failed to create bucket', 'error');
        }
    } catch (error) {
        showToast('Failed to create bucket', 'error');
    } finally {
        createBucketBtn.disabled = false;
    }
}

// Delete a bucket (only works if bucket is empty)
async function deleteBucket(bucketName) {
    if (!confirm(`Are you sure you want to delete bucket '${bucketName}'? This cannot be undone.`)) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE_URL}/buckets/${bucketName}`, {
            method: 'DELETE'
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showToast(`Bucket '${bucketName}' deleted successfully`, 'success');
            
            // If the current bucket was deleted, reset the view
            if (currentBucket === bucketName) {
                currentBucket = '';
                currentBucketHeader.textContent = 'Select a bucket';
                filesList.innerHTML = '<div class="no-files">No bucket selected</div>';
            }
            
            // Reload the bucket list
            loadBuckets();
        } else {
            showToast(data.error || `Failed to delete bucket '${bucketName}'`, 'error');
        }
    } catch (error) {
        showToast(`Failed to delete bucket '${bucketName}'`, 'error');
    }
}

// Show toast notification
function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    
    toastContainer.appendChild(toast);
    
    // Remove after 3 seconds
    setTimeout(() => {
        toast.remove();
    }, 3000);
}
