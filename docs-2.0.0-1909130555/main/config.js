function createConfig() {
    return {
        home: "about/01_overview",
        release: "2.0.0-1909130555",
        releases: [
            "2.0.0-1909130555"
        ],
        pathColors: {
            "*": "blue-grey"
        },
        theme: {
            primary: '#1976D2',
            secondary: '#424242',
            accent: '#82B1FF',
            error: '#FF5252',
            info: '#2196F3',
            success: '#4CAF50',
            warning: '#FFC107'
        },
        navTitle: 'Coherence Operator',
        navIcon: null,
        navLogo: 'images/oracle-coherence.svg'
    };
}

function createRoutes(){
    return [
        {
            path: '/about/01_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Coherence Operator documentation',
                keywords: 'oracle coherence, kubernetes, operator, documentation',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-01_overview', '/about/01_overview', {})
        },
        {
            path: '/about/02_concepts',
            meta: {
                h1: 'Coherence Operator Concepts',
                title: 'Coherence Operator Concepts',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-02_concepts', '/about/02_concepts', {})
        },
        {
            path: '/about/03_kubernetes',
            meta: {
                h1: 'Kubernetes on your Desktop',
                title: 'Kubernetes on your Desktop',
                description: 'Running Kubernetes on your desktop.',
                keywords: 'kubernetes',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-03_kubernetes', '/about/03_kubernetes', {})
        },
        {
            path: '/install/01_introduction',
            meta: {
                h1: 'Coherence Operator Installation',
                title: 'Coherence Operator Installation',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('install-01_introduction', '/install/01_introduction', {})
        },
        {
            path: '/install/02_prerequisites',
            meta: {
                h1: 'Prerequisites',
                title: 'Prerequisites',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('install-02_prerequisites', '/install/02_prerequisites', {})
        },
        {
            path: '/install/03_helm_install',
            meta: {
                h1: 'Installing With Helm',
                title: 'Installing With Helm',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('install-03_helm_install', '/install/03_helm_install', {})
        },
        {
            path: '/install/04_manual_install',
            meta: {
                h1: 'Installing Manually',
                title: 'Installing Manually',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('install-04_manual_install', '/install/04_manual_install', {})
        },
        {
            path: '/install/05_pre_release_versions',
            meta: {
                h1: 'Accessing Pre-Release Versions',
                title: 'Accessing Pre-Release Versions',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('install-05_pre_release_versions', '/install/05_pre_release_versions', {})
        },
        {
            path: '/clusters/01_introduction',
            meta: {
                h1: 'Create Coherence Clusters',
                title: 'Create Coherence Clusters',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-01_introduction', '/clusters/01_introduction', {})
        },
        {
            path: '/guides/01_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Coherence Operator guides',
                keywords: 'oracle coherence, kubernetes, operator, guides',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-01_overview', '/guides/01_overview', {})
        },
        {
            path: '/guides/02_quickstart',
            meta: {
                h1: 'Quick Start',
                title: 'Quick Start',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-02_quickstart', '/guides/02_quickstart', {})
        },
        {
            path: '/guides/03_management',
            meta: {
                h1: 'Management over ReST',
                title: 'Management over ReST',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-03_management', '/guides/03_management', {})
        },
        {
            path: '/guides/04_metrics',
            meta: {
                h1: 'Metrics',
                title: 'Metrics',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-04_metrics', '/guides/04_metrics', {})
        },
        {
            path: '/guides/05_logging',
            meta: {
                h1: 'Logging with ELK',
                title: 'Logging with ELK',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-05_logging', '/guides/05_logging', {})
        },
        {
            path: '/developer/01_introduction',
            meta: {
                h1: 'Coherence Operator',
                title: 'Coherence Operator',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-01_introduction', '/developer/01_introduction', {})
        },
        {
            path: '/', redirect: '/about/01_overview'
        },
        {
            path: '*', redirect: '/'
        }
    ];
}

function createNav(){
    return [
        { header: 'Core documentation' },
        {
            title: 'About',
            action: 'assistant',
            group: '/about',
            items: [
                { href: '/about/01_overview', title: 'Overview' },
                { href: '/about/02_concepts', title: 'Coherence Operator Concepts' },
                { href: '/about/03_kubernetes', title: 'Kubernetes on your Desktop' }
            ]
        },
        {
            title: 'Installing',
            action: 'settings',
            group: '/install',
            items: [
                { href: '/install/01_introduction', title: 'Coherence Operator Installation' },
                { href: '/install/02_prerequisites', title: 'Prerequisites' },
                { href: '/install/03_helm_install', title: 'Installing With Helm' },
                { href: '/install/04_manual_install', title: 'Installing Manually' },
                { href: '/install/05_pre_release_versions', title: 'Accessing Pre-Release Versions' }
            ]
        },
        {
            title: 'Coherence Clusters',
            action: 'widgets',
            group: '/clusters',
            items: [
                { href: '/clusters/01_introduction', title: 'Create Coherence Clusters' }
            ]
        },
        {
            title: 'Guides',
            action: 'explore',
            group: '/guides',
            items: [
                { href: '/guides/01_overview', title: 'Overview' },
                { href: '/guides/02_quickstart', title: 'Quick Start' },
                { href: '/guides/03_management', title: 'Management over ReST' },
                { href: '/guides/04_metrics', title: 'Metrics' },
                { href: '/guides/05_logging', title: 'Logging with ELK' }
            ]
        },
        {
            title: 'Developer documentation',
            action: 'build',
            group: '/developer',
            items: [
                { href: '/developer/01_introduction', title: 'Coherence Operator' }
            ]
        },
        { divider: true },
        { header: 'Additional resources' },
        {
            title: 'Community',
            action: 'fa-slack',
            href: 'https://join.slack.com/t/oraclecoherence/shared_invite/enQtNjA3MTU3MTk0MTE3LWZhMTdhM2E0ZDY2Y2FmZDhiOThlYzJjYTc5NzdkYWVlMzUzODZiNTI4ZWU3ZTlmNDQ4MmE1OTRhOWI1MmIxZjQ',
            target: '_blank'
        },
        {
            title: 'GitHub',
            action: 'fa-github-square',
            href: 'https://github.com/oracle/coherence-operator',
            target: '_blank'
        }
    ];
}