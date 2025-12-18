document.addEventListener('DOMContentLoaded', function() {
    const urlForm = document.getElementById('urlForm');
    const videoUrlInput = document.getElementById('videoUrl');
    const loadingSection = document.getElementById('loadingSection');
    const videoInfoSection = document.getElementById('videoInfoSection');
    const downloadProgress = document.getElementById('downloadProgress');
    const toggleLog = document.getElementById('toggleLog');
    const logContent = document.getElementById('logContent');
    const logText = document.getElementById('logText');
    const progressBar = document.getElementById('progressBar');
    const progressText = document.getElementById('progressText');
    const fileNameEl = document.getElementById('fileName');
    const fileSizeEl = document.getElementById('fileSize');
    const speedEl = document.getElementById('speed');
    const etaEl = document.getElementById('eta');

    let currentVideoInfo = null;
    let progressInterval = null;

    // Tab functionality
    const tabButtons = document.querySelectorAll('.tab-button');
    const formatLists = document.querySelectorAll('.format-list');

    tabButtons.forEach(button => {
        button.addEventListener('click', () => {
            const tab = button.getAttribute('data-tab');

            // Update active button
            tabButtons.forEach(btn => btn.classList.remove('active'));
            button.classList.add('active');

            // Show selected format list
            formatLists.forEach(list => {
                if (list.id === `${tab}Formats`) {
                    list.style.display = 'block';
                } else {
                    list.style.display = 'none';
                }
            });
        });
    });

    // Toggle log visibility
    toggleLog.addEventListener('click', () => {
        const logContent = document.getElementById('logContent');
        if (logContent.style.display === 'none') {
            logContent.style.display = 'block';
            toggleLog.textContent = '-';
            toggleLog.setAttribute('aria-expanded', 'true');
        } else {
            logContent.style.display = 'none';
            toggleLog.textContent = '+';
            toggleLog.setAttribute('aria-expanded', 'false');
        }
    });

    // Handle form submission
    urlForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const videoUrl = videoUrlInput.value.trim();
        if (!videoUrl) {
            showError('Silakan masukkan URL video');
            return;
        }

        // Show loading section
        loadingSection.style.display = 'block';
        videoInfoSection.style.display = 'none';
        downloadProgress.style.display = 'none';
        clearMessages();

        try {
            const response = await fetch('/info', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `url=${encodeURIComponent(videoUrl)}`
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            currentVideoInfo = await response.json();
            displayVideoInfo(currentVideoInfo);
            showSuccess('Informasi video berhasil dimuat');
        } catch (error) {
            console.error('Error:', error);
            showError('Terjadi kesalahan saat mengambil informasi video: ' + error.message);
        } finally {
            loadingSection.style.display = 'none';
        }
    });

    // Function to display video information
    function displayVideoInfo(videoInfo) {
        document.getElementById('videoThumbnail').src = videoInfo.thumbnails.length > 0 ?
            videoInfo.thumbnails[0].url : 'data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="150" viewBox="0 0 24 24"><rect width="24" height="18" fill="%23f0f0f0"/><text x="12" y="12" font-size="4" font-family="Arial" text-anchor="middle" fill="%23999">No Image</text></svg>';
        document.getElementById('videoTitle').textContent = videoInfo.title;
        document.getElementById('videoDuration').textContent = `Durasi: ${formatDuration(videoInfo.duration)}`;

        // Separate video and audio formats
        const videoFormats = [];
        const audioFormats = [];

        videoInfo.formats.forEach(format => {
            if (format.audio_only) {
                audioFormats.push(format);
            } else if (format.video_only) {
                videoFormats.push(format);
            } else if (format.resolution) {
                // Format yang memiliki video dan audio
                videoFormats.push(format);
            }
        });

        displayFormats(videoFormats, 'videoFormats');
        displayFormats(audioFormats, 'audioFormats');

        // Show video info section
        videoInfoSection.style.display = 'block';

        // Add event to download button
        const downloadBtn = document.getElementById('downloadBtn');
        downloadBtn.addEventListener('click', () => {
            startDownload(videoInfo);
        });

        // Reset format selection
        downloadBtn.dataset.selectedFormat = '';
    }

    // Function to display formats in the UI
    function displayFormats(formats, containerId) {
        const container = document.getElementById(containerId);
        container.innerHTML = '';

        if (formats.length === 0) {
            container.innerHTML = '<p class="text-muted">Tidak ada format tersedia</p>';
            return;
        }

        formats.forEach(format => {
            const formatElement = document.createElement('button');
            formatElement.type = 'button';
            formatElement.className = 'format-item p-2 mb-2 border rounded w-100 text-start';
            formatElement.setAttribute('aria-pressed', 'false');
            formatElement.innerHTML = `
                <div class="format-details d-flex justify-content-between align-items-center">
                    <span><strong>${format.format_id}</strong> - ${format.ext} - ${format.resolution || 'Audio only'}</span>
                    <span class="badge bg-secondary">${format.filesize ? formatFileSize(format.filesize) : 'Ukuran tidak diketahui'}</span>
                </div>
                <small class="text-muted">${format.format_note || 'No description'}</small>
            `;

            formatElement.addEventListener('click', () => {
                // Remove selection from other items
                document.querySelectorAll('#' + containerId + ' .format-item').forEach(item => {
                    item.classList.remove('border-primary', 'bg-light');
                    item.setAttribute('aria-pressed', 'false');
                });

                // Select this item
                formatElement.classList.add('border-primary', 'bg-light');
                formatElement.setAttribute('aria-pressed', 'true');

                // Store selected format in data attribute
                document.getElementById('downloadBtn').dataset.selectedFormat = format.format_id;
            });

            container.appendChild(formatElement);
        });
    }

    // Function to start download
    async function startDownload(videoInfo) {
        const selectedFormat = document.querySelector('#downloadBtn').dataset.selectedFormat;

        if (!selectedFormat) {
            showError('Silakan pilih format terlebih dahulu');
            return;
        }

        // Show download progress section
        downloadProgress.style.display = 'block';
        videoInfoSection.style.display = 'none';
        clearMessages();

        // Clear previous log
        logText.textContent = '';

        // Reset progress bar
        progressBar.style.width = '0%';
        progressText.textContent = '0%';
        fileNameEl.textContent = 'Nama file: Belum diketahui';
        fileSizeEl.textContent = 'Ukuran: Belum diketahui';
        speedEl.textContent = 'Kecepatan: Belum diketahui';
        etaEl.textContent = 'Perkiraan Selesai: Belum diketahui';

        // Create download request
        const downloadReq = {
            url: videoInfo.raw_response.webpage_url || window.location.href,
            format: selectedFormat,
            output_folder: './downloads',
            filename: '%(title)s.%(ext)s'
        };

        try {
            // Send the download request and handle streaming response
            const response = await fetch('/download', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(downloadReq)
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // Handle server-sent events for real-time logging
            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let buffer = '';

            while (true) {
                const { done, value } = await reader.read();

                if (done) {
                    break;
                }

                // Decode the chunk
                const chunk = decoder.decode(value, { stream: true });
                buffer += chunk;

                // Process each line (SSE format)
                const lines = buffer.split('\n');
                buffer = lines.pop(); // Keep incomplete line in buffer

                for (const line of lines) {
                    if (line.trim() === '') continue;

                    if (line.startsWith('data: ')) {
                        const data = line.slice(6); // Remove 'data: ' prefix

                        // Check if this is the finish signal
                        if (data === 'FINISHED') {
                            logText.textContent += '\n\nDownload selesai!';

                            // Parse the last response to get file info for auto-download
                            const fileInfo = extractFileInfo(logText.textContent);
                            if (fileInfo.fileName) {
                                // Create auto-download link
                                const downloadLink = `/downlink?reff=${encodeURIComponent(fileInfo.fileName)}`;
                                logText.textContent += `\n\nFile siap diunduh: ${downloadLink}`;

                                // Auto-hide progress and show completion message
                                setTimeout(() => {
                                    showSuccess('Download selesai! File sedang dipersiapkan untuk diunduh.');

                                    // Automatically open download link
                                    window.open(downloadLink, '_blank');
                                }, 1000);
                            }
                            break;
                        } else if (data.startsWith('TIMEOUT')) {
                            logText.textContent += `\n${data}`;
                            showError(data);
                            break;
                        } else {
                            // Append log data
                            logText.textContent += `\n${data}`;

                            // Parse progress information from yt-dlp output
                            parseProgress(data);

                            // Auto-scroll to bottom
                            logText.scrollTop = logText.scrollHeight;
                        }
                    }
                }
            }
        } catch (error) {
            console.error('Error:', error);
            logText.textContent += `\nError: ${error.message}`;
            showError('Terjadi kesalahan saat proses download: ' + error.message);
        }
    }

    // Parse progress information from yt-dlp log
    function parseProgress(logLine) {
        // Regular expressions to extract progress information
        const percentMatch = logLine.match(/(\d+\.\d+)%/);
        const speedMatch = logLine.match(/Speed:\s*([\d\.]+\s*[KMGT]iB\/s)/);
        const etaMatch = logLine.match(/ETA:\s*(\d+:\d+:\d+)/);
        const fileSizeMatch = logLine.match(/of\s*([0-9.]+\s*[KMGT]iB)(?:\s*\()/);
        const fileNameMatch = logLine.match(/\[download\]\s*(.+?)\s*has\s*already/);

        if (percentMatch) {
            const percent = parseFloat(percentMatch[1]);
            progressBar.style.width = `${percent}%`;
            progressText.textContent = `${percent}%`;
        }

        if (speedMatch) {
            speedEl.textContent = `Kecepatan: ${speedMatch[1]}`;
        }

        if (etaMatch) {
            etaEl.textContent = `Perkiraan Selesai: ${etaMatch[1]}`;
        }

        if (fileSizeMatch) {
            fileSizeEl.textContent = `Ukuran: ${fileSizeMatch[1]}`;
        }

        if (fileNameMatch) {
            fileNameEl.textContent = `Nama file: ${fileNameMatch[1]}`;
        }
    }

    // Extract file information from log
    function extractFileInfo(logText) {
        // Find the downloaded file name from the log
        const lines = logText.split('\n');
        for (let i = lines.length - 1; i >= 0; i--) {
            const line = lines[i];

            // Look for a line containing the saved filename
            const matches = line.match(/Writing video description metadata as JSON to:\s*(.*)/);
            if (matches) {
                const filePath = matches[1];
                const fileName = filePath.split('/').pop().split('\\').pop();
                return { fileName };
            }

            // Alternative pattern for filename extraction
            const dlCompleteMatches = line.match(/Destination:\s*(.*)/);
            if (dlCompleteMatches) {
                const fileName = dlCompleteMatches[1].split('/').pop().split('\\').pop();
                return { fileName };
            }
        }

        return { fileName: null };
    }

    // Helper functions
    function formatDuration(seconds) {
        if (!seconds) return 'Durasi tidak diketahui';

        const hrs = Math.floor(seconds / 3600);
        const mins = Math.floor((seconds % 3600) / 60);
        const secs = Math.floor(seconds % 60);

        if (hrs > 0) {
            return `${hrs}:${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
        }
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    }

    function formatFileSize(bytes) {
        if (!bytes) return 'Ukuran tidak diketahui';

        if (bytes < 1024) return bytes + ' B';
        else if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' KB';
        else if (bytes < 1073741824) return (bytes / 1048576).toFixed(1) + ' MB';
        else return (bytes / 1073741824).toFixed(1) + ' GB';
    }

    function showError(message) {
        // Membuat elemen alert error Bootstrap
        let errorAlert = document.querySelector('.alert-danger');
        if (!errorAlert) {
            errorAlert = document.createElement('div');
            errorAlert.className = 'alert alert-danger alert-dismissible fade show mt-3';
            errorAlert.setAttribute('role', 'alert');
            errorAlert.innerHTML = `
                <span class="error-message"></span>
                <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
            `;
            document.querySelector('.container').insertBefore(errorAlert, document.querySelector('header').nextSibling);
        }
        errorAlert.querySelector('.error-message').textContent = message;

        // Sembunyikan pesan sukses
        const successAlert = document.querySelector('.alert-success');
        if (successAlert) {
            successAlert.style.display = 'none';
        }
    }

    function showSuccess(message) {
        // Membuat elemen alert sukses Bootstrap
        let successAlert = document.querySelector('.alert-success');
        if (!successAlert) {
            successAlert = document.createElement('div');
            successAlert.className = 'alert alert-success alert-dismissible fade show mt-3';
            successAlert.setAttribute('role', 'alert');
            successAlert.innerHTML = `
                <span class="success-message"></span>
                <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
            `;
            document.querySelector('.container').insertBefore(successAlert, document.querySelector('header').nextSibling);
        }
        successAlert.querySelector('.success-message').textContent = message;
        successAlert.style.display = 'block';

        // Sembunyikan pesan error
        const errorAlert = document.querySelector('.alert-danger');
        if (errorAlert) {
            errorAlert.style.display = 'none';
        }
    }

    function clearMessages() {
        const errorAlert = document.querySelector('.alert-danger');
        const successAlert = document.querySelector('.alert-success');
        if (errorAlert) errorAlert.style.display = 'none';
        if (successAlert) successAlert.style.display = 'none';
    }
});