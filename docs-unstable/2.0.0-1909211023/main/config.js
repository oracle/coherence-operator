function createConfig() {
    return {
        home: "about/01_overview",
        release: "2.0.0-1909211023",
        releases: [
            "2.0.0-1909211023"
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
            path: '/clusters/010_introduction',
            meta: {
                h1: 'Create Coherence Clusters',
                title: 'Create Coherence Clusters',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-010_introduction', '/clusters/010_introduction', {})
        },
        {
            path: '/clusters/020_k8s_resources',
            meta: {
                h1: 'Coherence Clusters K8s Resources',
                title: 'Coherence Clusters K8s Resources',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-020_k8s_resources', '/clusters/020_k8s_resources', {})
        },
        {
            path: '/clusters/030_roles',
            meta: {
                h1: 'Define Coherence Roles',
                title: 'Define Coherence Roles',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-030_roles', '/clusters/030_roles', {})
        },
        {
            path: '/clusters/040_replicas',
            meta: {
                h1: 'Role Replica Count',
                title: 'Role Replica Count',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-040_replicas', '/clusters/040_replicas', {})
        },
        {
            path: '/clusters/045_storage_enabled',
            meta: {
                h1: 'Storage Enabled or Disabled Roles',
                title: 'Storage Enabled or Disabled Roles',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-045_storage_enabled', '/clusters/045_storage_enabled', {})
        },
        {
            path: '/clusters/050_config_files',
            meta: {
                h1: 'Coherence Config Files',
                title: 'Coherence Config Files',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-050_config_files', '/clusters/050_config_files', {})
        },
        {
            path: '/clusters/060_coherence_image',
            meta: {
                h1: 'Setting the Coherence Image',
                title: 'Setting the Coherence Image',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-060_coherence_image', '/clusters/060_coherence_image', {})
        },
        {
            path: '/clusters/065_application_image',
            meta: {
                h1: 'Setting the Application Image',
                title: 'Setting the Application Image',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-065_application_image', '/clusters/065_application_image', {})
        },
        {
            path: '/clusters/070_private_repos',
            meta: {
                h1: 'Using Private Image Registries',
                title: 'Using Private Image Registries',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-070_private_repos', '/clusters/070_private_repos', {})
        },
        {
            path: '/clusters/080_environment_variables',
            meta: {
                h1: 'Environment Variables',
                title: 'Environment Variables',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-080_environment_variables', '/clusters/080_environment_variables', {})
        },
        {
            path: '/clusters/100_logging',
            meta: {
                h1: 'Logging Configuration',
                title: 'Logging Configuration',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-100_logging', '/clusters/100_logging', {})
        },
        {
            path: '/clusters/120_annotations',
            meta: {
                h1: 'Pod Annotations',
                title: 'Pod Annotations',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-120_annotations', '/clusters/120_annotations', {})
        },
        {
            path: '/clusters/150_volumes',
            meta: {
                h1: 'Configure Volumes',
                title: 'Configure Volumes',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-150_volumes', '/clusters/150_volumes', {})
        },
        {
            path: '/guides/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Coherence Operator guides',
                keywords: 'oracle coherence, kubernetes, operator, guides',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-010_overview', '/guides/010_overview', {})
        },
        {
            path: '/guides/020_quickstart',
            meta: {
                h1: 'Quick Start',
                title: 'Quick Start',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-020_quickstart', '/guides/020_quickstart', {})
        },
        {
            path: '/guides/030_applications',
            meta: {
                h1: 'Deploying Coherence Applications',
                title: 'Deploying Coherence Applications',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-030_applications', '/guides/030_applications', {})
        },
        {
            path: '/guides/080_management',
            meta: {
                h1: 'Management over ReST',
                title: 'Management over ReST',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-080_management', '/guides/080_management', {})
        },
        {
            path: '/guides/090_metrics',
            meta: {
                h1: 'Metrics',
                title: 'Metrics',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-090_metrics', '/guides/090_metrics', {})
        },
        {
            path: '/guides/100_logging',
            meta: {
                h1: 'Logging with ELK',
                title: 'Logging with ELK',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('guides-100_logging', '/guides/100_logging', {})
        },
        {
            path: '/developer/01_introduction',
            meta: {
                h1: 'Coherence Operator Development',
                title: 'Coherence Operator Development',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-01_introduction', '/developer/01_introduction', {})
        },
        {
            path: '/developer/03_high_level',
            meta: {
                h1: 'High Level Design',
                title: 'High Level Design',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-03_high_level', '/developer/03_high_level', {})
        },
        {
            path: '/developer/04_how_it_works',
            meta: {
                h1: 'How It Works',
                title: 'How It Works',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-04_how_it_works', '/developer/04_how_it_works', {})
        },
        {
            path: '/developer/05_building',
            meta: {
                h1: 'Building and Testing',
                title: 'Building and Testing',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-05_building', '/developer/05_building', {})
        },
        {
            path: '/developer/06_debugging',
            meta: {
                h1: 'Debugging',
                title: 'Debugging',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-06_debugging', '/developer/06_debugging', {})
        },
        {
            path: '/developer/07_execution',
            meta: {
                h1: 'Execution',
                title: 'Execution',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-07_execution', '/developer/07_execution', {})
        },
        {
            path: '/developer/08_docs',
            meta: {
                h1: 'Building the Docs',
                title: 'Building the Docs',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-08_docs', '/developer/08_docs', {})
        },
        {
            path: '/developer/09_useful',
            meta: {
                h1: 'Useful Info',
                title: 'Useful Info',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('developer-09_useful', '/developer/09_useful', {})
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
                { href: '/clusters/010_introduction', title: 'Create Coherence Clusters' },
                { href: '/clusters/020_k8s_resources', title: 'Coherence Clusters K8s Resources' },
                { href: '/clusters/030_roles', title: 'Define Coherence Roles' },
                { href: '/clusters/040_replicas', title: 'Role Replica Count' },
                { href: '/clusters/045_storage_enabled', title: 'Storage Enabled or Disabled Roles' },
                { href: '/clusters/050_config_files', title: 'Coherence Config Files' },
                { href: '/clusters/060_coherence_image', title: 'Setting the Coherence Image' },
                { href: '/clusters/065_application_image', title: 'Setting the Application Image' },
                { href: '/clusters/070_private_repos', title: 'Using Private Image Registries' },
                { href: '/clusters/080_environment_variables', title: 'Environment Variables' },
                { href: '/clusters/100_logging', title: 'Logging Configuration' },
                { href: '/clusters/120_annotations', title: 'Pod Annotations' },
                { href: '/clusters/150_volumes', title: 'Configure Volumes' }
            ]
        },
        {
            title: 'Guides',
            action: 'explore',
            group: '/guides',
            items: [
                { href: '/guides/010_overview', title: 'Overview' },
                { href: '/guides/020_quickstart', title: 'Quick Start' },
                { href: '/guides/030_applications', title: 'Deploying Coherence Applications' },
                { href: '/guides/080_management', title: 'Management over ReST' },
                { href: '/guides/090_metrics', title: 'Metrics' },
                { href: '/guides/100_logging', title: 'Logging with ELK' }
            ]
        },
        {
            title: 'Developer documentation',
            action: 'build',
            group: '/developer',
            items: [
                { href: '/developer/01_introduction', title: 'Coherence Operator Development' },
                { href: '/developer/03_high_level', title: 'High Level Design' },
                { href: '/developer/04_how_it_works', title: 'How It Works' },
                { href: '/developer/05_building', title: 'Building and Testing' },
                { href: '/developer/06_debugging', title: 'Debugging' },
                { href: '/developer/07_execution', title: 'Execution' },
                { href: '/developer/08_docs', title: 'Building the Docs' },
                { href: '/developer/09_useful', title: 'Useful Info' }
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