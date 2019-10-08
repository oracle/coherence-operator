function createConfig() {
    return {
        home: "about/01_overview",
        release: "2.0.0-1910081143",
        releases: [
            "2.0.0-1910081143"
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
            path: '/about/04_quickstart',
            meta: {
                h1: 'Quick Start',
                title: 'Quick Start',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-04_quickstart', '/about/04_quickstart', {})
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
            path: '/clusters/045_coherence_storage_enabled',
            meta: {
                h1: 'Storage Enabled or Disabled Roles',
                title: 'Storage Enabled or Disabled Roles',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-045_coherence_storage_enabled', '/clusters/045_coherence_storage_enabled', {})
        },
        {
            path: '/clusters/050_coherence_config_files',
            meta: {
                h1: 'Coherence Config Files',
                title: 'Coherence Config Files',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-050_coherence_config_files', '/clusters/050_coherence_config_files', {})
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
            path: '/clusters/070_applications',
            meta: {
                h1: 'Configure Applications',
                title: 'Configure Applications',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-070_applications', '/clusters/070_applications', {})
        },
        {
            path: '/clusters/080_jvm_heap_size',
            meta: {
                h1: 'Setting the JVM Heap Size',
                title: 'Setting the JVM Heap Size',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-080_jvm_heap_size', '/clusters/080_jvm_heap_size', {})
        },
        {
            path: '/clusters/090_environment_variables',
            meta: {
                h1: 'Environment Variables',
                title: 'Environment Variables',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-090_environment_variables', '/clusters/090_environment_variables', {})
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
                h1: 'Configure Pod Annotations',
                title: 'Configure Pod Annotations',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-120_annotations', '/clusters/120_annotations', {})
        },
        {
            path: '/clusters/125_labels',
            meta: {
                h1: 'Configure Pod Labels',
                title: 'Configure Pod Labels',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-125_labels', '/clusters/125_labels', {})
        },
        {
            path: '/clusters/150_volumes',
            meta: {
                h1: 'Configure Additional Volumes',
                title: 'Configure Additional Volumes',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-150_volumes', '/clusters/150_volumes', {})
        },
        {
            path: '/clusters/190_service_account',
            meta: {
                h1: 'Kubernetes Service Account',
                title: 'Kubernetes Service Account',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-190_service_account', '/clusters/190_service_account', {})
        },
        {
            path: '/clusters/200_private_repos',
            meta: {
                h1: 'Using Private Image Registries',
                title: 'Using Private Image Registries',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-200_private_repos', '/clusters/200_private_repos', {})
        },
        {
            path: '/app-deployments/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Application Deployments',
                keywords: 'oracle coherence, kubernetes, operator, Application Deployments',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-010_overview', '/app-deployments/010_overview', {})
        },
        {
            path: '/app-deployments/020_packaging',
            meta: {
                h1: 'Packaging Applications',
                title: 'Packaging Applications',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-020_packaging', '/app-deployments/020_packaging', {})
        },
        {
            path: '/app-deployments/030_extend',
            meta: {
                h1: 'Using Coherence*Extend',
                title: 'Using Coherence*Extend',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-030_extend', '/app-deployments/030_extend', {})
        },
        {
            path: '/app-deployments/040_storage_disabled',
            meta: {
                h1: 'Storage-Disabled Clients',
                title: 'Storage-Disabled Clients',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-040_storage_disabled', '/app-deployments/040_storage_disabled', {})
        },
        {
            path: '/app-deployments/050_federation',
            meta: {
                h1: 'Federation',
                title: 'Federation',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-050_federation', '/app-deployments/050_federation', {})
        },
        {
            path: '/app-deployments/060_persistence',
            meta: {
                h1: 'Persistence',
                title: 'Persistence',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-060_persistence', '/app-deployments/060_persistence', {})
        },
        {
            path: '/app-deployments/070_elasticdata',
            meta: {
                h1: 'Elastic Data',
                title: 'Elastic Data',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-070_elasticdata', '/app-deployments/070_elasticdata', {})
        },
        {
            path: '/app-deployments/080_scaling',
            meta: {
                h1: 'Scaling Deployments',
                title: 'Scaling Deployments',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-080_scaling', '/app-deployments/080_scaling', {})
        },
        {
            path: '/app-deployments/090_rolling',
            meta: {
                h1: 'Rolling Upgrades',
                title: 'Rolling Upgrades',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployments-090_rolling', '/app-deployments/090_rolling', {})
        },
        {
            path: '/examples/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Coherence Operator examples',
                keywords: 'oracle coherence, kubernetes, operator, examples',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-010_overview', '/examples/010_overview', {})
        },
        {
            path: '/examples/011_prereqs',
            meta: {
                h1: 'Example Pre-requisites',
                title: 'Example Pre-requisites',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-011_prereqs', '/examples/011_prereqs', {})
        },
        {
            path: '/examples/015_packaging',
            meta: {
                h1: 'Packaging Applications',
                title: 'Packaging Applications',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-015_packaging', '/examples/015_packaging', {})
        },
        {
            path: '/examples/020_extend',
            meta: {
                h1: 'Coherence*Extend Deployments',
                title: 'Coherence*Extend Deployments',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-020_extend', '/examples/020_extend', {})
        },
        {
            path: '/examples/030_storage_disabled',
            meta: {
                h1: 'Storage-Disabled Deployments',
                title: 'Storage-Disabled Deployments',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-030_storage_disabled', '/examples/030_storage_disabled', {})
        },
        {
            path: '/examples/045_rolling',
            meta: {
                h1: 'Rolling Upgrades',
                title: 'Rolling Upgrades',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-045_rolling', '/examples/045_rolling', {})
        },
        {
            path: '/examples/046_federation',
            meta: {
                h1: 'Federation',
                title: 'Federation',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-046_federation', '/examples/046_federation', {})
        },
        {
            path: '/examples/050_custom',
            meta: {
                h1: 'Configure a Custom Logger',
                title: 'Configure a Custom Logger',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-050_custom', '/examples/050_custom', {})
        },
        {
            path: '/examples/060_ssl',
            meta: {
                h1: 'Configuring a Custom Logger',
                title: 'Configuring a Custom Logger',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-060_ssl', '/examples/060_ssl', {})
        },
        {
            path: '/examples/070_elastic1',
            meta: {
                h1: 'Elastic Data - Defaults',
                title: 'Elastic Data - Defaults',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-070_elastic1', '/examples/070_elastic1', {})
        },
        {
            path: '/examples/080_elastic2',
            meta: {
                h1: 'Elastic Data - Custom',
                title: 'Elastic Data - Custom',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('examples-080_elastic2', '/examples/080_elastic2', {})
        },
        {
            path: '/metrics/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Metrics',
                keywords: 'oracle coherence, kubernetes, operator, Metrics',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-010_overview', '/metrics/010_overview', {})
        },
        {
            path: '/metrics/020_metrics',
            meta: {
                h1: 'Enabling Metrics',
                title: 'Enabling Metrics',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-020_metrics', '/metrics/020_metrics', {})
        },
        {
            path: '/metrics/030_ssl',
            meta: {
                h1: 'Enabling SSL',
                title: 'Enabling SSL',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-030_ssl', '/metrics/030_ssl', {})
        },
        {
            path: '/metrics/040_scrapingl',
            meta: {
                h1: 'Using Your Own Prometheus',
                title: 'Using Your Own Prometheus',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-040_scrapingl', '/metrics/040_scrapingl', {})
        },
        {
            path: '/logging/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Logging',
                keywords: 'oracle coherence, kubernetes, operator, Logging',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-010_overview', '/logging/010_overview', {})
        },
        {
            path: '/logging/020_logging',
            meta: {
                h1: 'Enabling Log Capture',
                title: 'Enabling Log Capture',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-020_logging', '/logging/020_logging', {})
        },
        {
            path: '/logging/030_custom',
            meta: {
                h1: 'Using a Custom Logger',
                title: 'Using a Custom Logger',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-030_custom', '/logging/030_custom', {})
        },
        {
            path: '/logging/040_own',
            meta: {
                h1: 'Using Your Own Elasticsearch',
                title: 'Using Your Own Elasticsearch',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-040_own', '/logging/040_own', {})
        },
        {
            path: '/diagnostics/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Diagnostic Tools',
                keywords: 'oracle coherence, kubernetes, operator, Diagnostic Tools',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('diagnostics-010_overview', '/diagnostics/010_overview', {})
        },
        {
            path: '/diagnostics/020_heapdump',
            meta: {
                h1: 'Generating Heap Dumps',
                title: 'Generating Heap Dumps',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('diagnostics-020_heapdump', '/diagnostics/020_heapdump', {})
        },
        {
            path: '/diagnostics/030_jfr',
            meta: {
                h1: 'Using Java Flight Recorder',
                title: 'Using Java Flight Recorder',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('diagnostics-030_jfr', '/diagnostics/030_jfr', {})
        },
        {
            path: '/diagnostics/040_reporter',
            meta: {
                h1: 'Coherence Reporter',
                title: 'Coherence Reporter',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('diagnostics-040_reporter', '/diagnostics/040_reporter', {})
        },
        {
            path: '/diagnostics/050_console',
            meta: {
                h1: 'Accessing the Coherence Console',
                title: 'Accessing the Coherence Console',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('diagnostics-050_console', '/diagnostics/050_console', {})
        },
        {
            path: '/diagnostics/060_cohql',
            meta: {
                h1: 'Accessing CohQL',
                title: 'Accessing CohQL',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('diagnostics-060_cohql', '/diagnostics/060_cohql', {})
        },
        {
            path: '/management/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Management Over ReST',
                keywords: 'oracle coherence, kubernetes, operator, Management, ReST',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-010_overview', '/management/010_overview', {})
        },
        {
            path: '/management/020_enabling',
            meta: {
                h1: 'Enabling',
                title: 'Enabling',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-020_enabling', '/management/020_enabling', {})
        },
        {
            path: '/management/030_visualvm',
            meta: {
                h1: 'Using JVisualVM',
                title: 'Using JVisualVM',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-030_visualvm', '/management/030_visualvm', {})
        },
        {
            path: '/management/040_ssl',
            meta: {
                h1: 'Enabling SSL',
                title: 'Enabling SSL',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-040_ssl', '/management/040_ssl', {})
        },
        {
            path: '/management/050_mbeans',
            meta: {
                h1: 'Modifying MBeans',
                title: 'Modifying MBeans',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-050_mbeans', '/management/050_mbeans', {})
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
                { href: '/about/03_kubernetes', title: 'Kubernetes on your Desktop' },
                { href: '/about/04_quickstart', title: 'Quick Start' }
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
                { href: '/clusters/045_coherence_storage_enabled', title: 'Storage Enabled or Disabled Roles' },
                { href: '/clusters/050_coherence_config_files', title: 'Coherence Config Files' },
                { href: '/clusters/060_coherence_image', title: 'Setting the Coherence Image' },
                { href: '/clusters/070_applications', title: 'Configure Applications' },
                { href: '/clusters/080_jvm_heap_size', title: 'Setting the JVM Heap Size' },
                { href: '/clusters/090_environment_variables', title: 'Environment Variables' },
                { href: '/clusters/100_logging', title: 'Logging Configuration' },
                { href: '/clusters/120_annotations', title: 'Configure Pod Annotations' },
                { href: '/clusters/125_labels', title: 'Configure Pod Labels' },
                { href: '/clusters/150_volumes', title: 'Configure Additional Volumes' },
                { href: '/clusters/190_service_account', title: 'Kubernetes Service Account' },
                { href: '/clusters/200_private_repos', title: 'Using Private Image Registries' }
            ]
        },
        {
            title: 'Application Deployments',
            action: 'extension',
            group: '/app-deployments',
            items: [
                { href: '/app-deployments/010_overview', title: 'Overview' },
                { href: '/app-deployments/020_packaging', title: 'Packaging Applications' },
                { href: '/app-deployments/030_extend', title: 'Using Coherence*Extend' },
                { href: '/app-deployments/040_storage_disabled', title: 'Storage-Disabled Clients' },
                { href: '/app-deployments/050_federation', title: 'Federation' },
                { href: '/app-deployments/060_persistence', title: 'Persistence' },
                { href: '/app-deployments/070_elasticdata', title: 'Elastic Data' },
                { href: '/app-deployments/080_scaling', title: 'Scaling Deployments' },
                { href: '/app-deployments/090_rolling', title: 'Rolling Upgrades' }
            ]
        },
        {
            title: 'Examples',
            action: 'explore',
            group: '/examples',
            items: [
                { href: '/examples/010_overview', title: 'Overview' },
                { href: '/examples/011_prereqs', title: 'Example Pre-requisites' },
                { href: '/examples/015_packaging', title: 'Packaging Applications' },
                { href: '/examples/020_extend', title: 'Coherence*Extend Deployments' },
                { href: '/examples/030_storage_disabled', title: 'Storage-Disabled Deployments' },
                { href: '/examples/045_rolling', title: 'Rolling Upgrades' },
                { href: '/examples/046_federation', title: 'Federation' },
                { href: '/examples/050_custom', title: 'Configure a Custom Logger' },
                { href: '/examples/060_ssl', title: 'Configuring a Custom Logger' },
                { href: '/examples/070_elastic1', title: 'Elastic Data - Defaults' },
                { href: '/examples/080_elastic2', title: 'Elastic Data - Custom' }
            ]
        },
        {
            title: 'Metrics',
            action: 'av_timer',
            group: '/metrics',
            items: [
                { href: '/metrics/010_overview', title: 'Overview' },
                { href: '/metrics/020_metrics', title: 'Enabling Metrics' },
                { href: '/metrics/030_ssl', title: 'Enabling SSL' },
                { href: '/metrics/040_scrapingl', title: 'Using Your Own Prometheus' }
            ]
        },
        {
            title: 'Logging',
            action: 'donut_large',
            group: '/logging',
            items: [
                { href: '/logging/010_overview', title: 'Overview' },
                { href: '/logging/020_logging', title: 'Enabling Log Capture' },
                { href: '/logging/030_custom', title: 'Using a Custom Logger' },
                { href: '/logging/040_own', title: 'Using Your Own Elasticsearch' }
            ]
        },
        {
            title: 'Diagnostic Tools',
            action: 'favorite_outline',
            group: '/diagnostics',
            items: [
                { href: '/diagnostics/010_overview', title: 'Overview' },
                { href: '/diagnostics/020_heapdump', title: 'Generating Heap Dumps' },
                { href: '/diagnostics/030_jfr', title: 'Using Java Flight Recorder' },
                { href: '/diagnostics/040_reporter', title: 'Coherence Reporter' },
                { href: '/diagnostics/050_console', title: 'Accessing the Coherence Console' },
                { href: '/diagnostics/060_cohql', title: 'Accessing CohQL' }
            ]
        },
        {
            title: 'Management over ReST',
            action: 'cloud',
            group: '/management',
            items: [
                { href: '/management/010_overview', title: 'Overview' },
                { href: '/management/020_enabling', title: 'Enabling' },
                { href: '/management/030_visualvm', title: 'Using JVisualVM' },
                { href: '/management/040_ssl', title: 'Enabling SSL' },
                { href: '/management/050_mbeans', title: 'Modifying MBeans' }
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