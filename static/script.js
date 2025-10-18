document.addEventListener('DOMContentLoaded', function() {
    const shortenForm = document.getElementById('shortenForm');
    const urlInput = document.getElementById('urlInput');
    const resultContainer = document.getElementById('resultContainer');
    const shortUrlResult = document.getElementById('shortUrlResult');
    const originalUrl = document.getElementById('originalUrl');
    const copyButton = document.getElementById('copyButton');
    const topDomains = document.getElementById('topDomains');

    // Load metrics on page load
    loadMetrics();

    // Handle form submission
    shortenForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const url = urlInput.value.trim();
        if (!url) return;

        shortenUrl(url);
    });

    // Handle copy button
    copyButton.addEventListener('click', function() {
        shortUrlResult.select();
        document.execCommand('copy');
        
        // Show feedback
        copyButton.textContent = 'Copied!';
        setTimeout(() => {
            copyButton.textContent = 'Copy';
        }, 2000);
    });

    // Function to shorten URL
    function shortenUrl(url) {
        fetch('/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ url: url })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            // Display result
            shortUrlResult.value = data.short_url;
            originalUrl.textContent = data.original_url;
            resultContainer.style.display = 'block';
            
            // Reload metrics
            loadMetrics();
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to shorten URL. Please try again.');
        });
    }

    // Function to load metrics
    function loadMetrics() {
        fetch('/metrics')
        .then(response => response.json())
        .then(data => {
            if (data.top_domains && data.top_domains.length > 0) {
                topDomains.innerHTML = '';
                
                data.top_domains.forEach(domain => {
                    const domainItem = document.createElement('div');
                    domainItem.className = 'domain-item';
                    
                    domainItem.innerHTML = `
                        <span class="domain-name">${domain.domain}</span>
                        <span class="domain-count">${domain.count}</span>
                    `;
                    
                    topDomains.appendChild(domainItem);
                });
            } else {
                topDomains.innerHTML = '<p>No domains shortened yet</p>';
            }
        })
        .catch(error => {
            console.error('Error loading metrics:', error);
            topDomains.innerHTML = '<p>Failed to load metrics</p>';
        });
    }
});