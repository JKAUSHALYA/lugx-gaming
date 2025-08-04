class LugxAnalytics {
    constructor(config = {}) {
        // Auto-detect environment and set appropriate API URL
        this.apiUrl = config.apiUrl || this.getApiUrl();
        this.sessionId = this.generateSessionId();
        this.startTime = new Date();
        this.pageStartTime = new Date();
        this.isActive = true;
        this.activeTime = 0;
        this.lastActiveTime = new Date();
        this.maxScrollPercentage = 0;
        this.clickCount = 0;
        this.pagesVisited = 1;
        
        this.init();
    }

    getApiUrl() {
        // Check if running in Kubernetes environment
        const hostname = window.location.hostname;
        if (hostname === 'localhost' || hostname === '127.0.0.1') {
            // Local development - try analytics service on different ports
            return 'http://localhost:30082/api/analytics';
        } else {
            // Production or cluster environment
            return `${window.location.protocol}//${hostname}:30082/api/analytics`;
        }
    }

    generateSessionId() {
        // Check if session ID exists in sessionStorage
        let sessionId = sessionStorage.getItem('lugx_session_id');
        if (!sessionId) {
            sessionId = 'sess_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
            sessionStorage.setItem('lugx_session_id', sessionId);
            sessionStorage.setItem('lugx_session_start', new Date().toISOString());
        }
        return sessionId;
    }

    init() {
        this.trackPageView();
        this.setupEventListeners();
        this.startActiveTimeTracking();
        this.setupBeforeUnloadHandler();
    }

    async sendData(endpoint, data) {
        try {
            await fetch(`${this.apiUrl}/${endpoint}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });
        } catch (error) {
            console.error('Analytics tracking error:', error);
        }
    }

    trackPageView() {
        const data = {
            session_id: this.sessionId,
            user_agent: navigator.userAgent,
            page_url: window.location.href,
            page_title: document.title,
            referrer: document.referrer,
            page_load_time: performance.timing ? 
                (performance.timing.loadEventEnd - performance.timing.navigationStart) : 0,
            viewport_width: window.innerWidth,
            viewport_height: window.innerHeight
        };

        this.sendData('pageview', data);
    }

    setupEventListeners() {
        // Track clicks
        document.addEventListener('click', (event) => {
            this.trackClick(event);
        });

        // Track scroll
        let scrollTimeout;
        window.addEventListener('scroll', () => {
            clearTimeout(scrollTimeout);
            scrollTimeout = setTimeout(() => {
                this.trackScrollDepth();
            }, 150);
        });

        // Track user activity
        ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart', 'click'].forEach(event => {
            document.addEventListener(event, () => {
                this.updateActiveTime();
            }, true);
        });

        // Track window focus/blur
        window.addEventListener('focus', () => {
            this.isActive = true;
            this.lastActiveTime = new Date();
        });

        window.addEventListener('blur', () => {
            this.isActive = false;
            this.updateActiveTime();
        });

        // Track page visibility changes
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                this.isActive = false;
                this.updateActiveTime();
            } else {
                this.isActive = true;
                this.lastActiveTime = new Date();
            }
        });
    }

    trackClick(event) {
        const element = event.target;
        const rect = element.getBoundingClientRect();
        
        const data = {
            session_id: this.sessionId,
            page_url: window.location.href,
            element_tag: element.tagName.toLowerCase(),
            element_id: element.id || '',
            element_class: element.className || '',
            element_text: element.textContent ? element.textContent.slice(0, 100) : '',
            click_x: Math.round(event.clientX),
            click_y: Math.round(event.clientY)
        };

        this.clickCount++;
        this.sendData('click', data);
    }

    trackScrollDepth() {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const documentHeight = Math.max(
            document.body.scrollHeight,
            document.body.offsetHeight,
            document.documentElement.clientHeight,
            document.documentElement.scrollHeight,
            document.documentElement.offsetHeight
        );
        const windowHeight = window.innerHeight;
        const scrollPercentage = Math.round((scrollTop + windowHeight) / documentHeight * 100);

        if (scrollPercentage > this.maxScrollPercentage) {
            this.maxScrollPercentage = Math.min(scrollPercentage, 100);
            
            const data = {
                session_id: this.sessionId,
                page_url: window.location.href,
                max_scroll_percentage: this.maxScrollPercentage,
                total_page_height: documentHeight,
                viewport_height: windowHeight
            };

            this.sendData('scroll', data);
        }
    }

    updateActiveTime() {
        if (this.isActive) {
            const now = new Date();
            this.activeTime += Math.round((now - this.lastActiveTime) / 1000);
            this.lastActiveTime = now;
        }
    }

    startActiveTimeTracking() {
        setInterval(() => {
            this.updateActiveTime();
        }, 1000);
    }

    trackPageTime() {
        const timeOnPage = Math.round((new Date() - this.pageStartTime) / 1000);
        
        const data = {
            session_id: this.sessionId,
            page_url: window.location.href,
            time_on_page: timeOnPage,
            is_active_time: this.activeTime
        };

        this.sendData('pagetime', data);
    }

    trackSessionTime() {
        const sessionStart = sessionStorage.getItem('lugx_session_start');
        const endTime = new Date();
        const totalDuration = Math.round((endTime - new Date(sessionStart)) / 1000);

        const data = {
            session_id: this.sessionId,
            start_time: sessionStart,
            end_time: endTime.toISOString(),
            total_session_duration: totalDuration,
            pages_visited: this.pagesVisited,
            total_clicks: this.clickCount,
            device_type: this.getDeviceType(),
            browser: this.getBrowser(),
            operating_system: this.getOperatingSystem()
        };

        this.sendData('sessiontime', data);
    }

    setupBeforeUnloadHandler() {
        window.addEventListener('beforeunload', () => {
            this.trackPageTime();
            this.trackSessionTime();
        });

        // Also track when page becomes hidden (for mobile browsers)
        document.addEventListener('visibilitychange', () => {
            if (document.hidden) {
                this.trackPageTime();
                this.trackSessionTime();
            }
        });
    }

    getDeviceType() {
        const userAgent = navigator.userAgent.toLowerCase();
        if (/tablet|ipad|playbook|silk/i.test(userAgent)) {
            return 'tablet';
        }
        if (/mobile|iphone|ipod|android|blackberry|opera|mini|windows\sce|palm|smartphone|iemobile/i.test(userAgent)) {
            return 'mobile';
        }
        return 'desktop';
    }

    getBrowser() {
        const userAgent = navigator.userAgent;
        if (userAgent.includes('Chrome')) return 'Chrome';
        if (userAgent.includes('Firefox')) return 'Firefox';
        if (userAgent.includes('Safari')) return 'Safari';
        if (userAgent.includes('Edge')) return 'Edge';
        if (userAgent.includes('Opera')) return 'Opera';
        return 'Unknown';
    }

    getOperatingSystem() {
        const userAgent = navigator.userAgent;
        if (userAgent.includes('Windows')) return 'Windows';
        if (userAgent.includes('Mac')) return 'macOS';
        if (userAgent.includes('Linux')) return 'Linux';
        if (userAgent.includes('Android')) return 'Android';
        if (userAgent.includes('iOS')) return 'iOS';
        return 'Unknown';
    }

    // Public method to manually track custom events
    trackCustomEvent(eventName, data = {}) {
        const customData = {
            session_id: this.sessionId,
            page_url: window.location.href,
            event_name: eventName,
            ...data
        };

        this.sendData('custom', customData);
    }
}

// Auto-initialize analytics when the script loads
document.addEventListener('DOMContentLoaded', () => {
    // Initialize with default configuration - apiUrl will be auto-detected
    window.lugxAnalytics = new LugxAnalytics();
});

// Also support manual initialization
window.LugxAnalytics = LugxAnalytics;
