function createConfig() {
    return {
        home: "about/01_overview",
        release: "2.0.1-1911091123",
        releases: [
            "2.0.1-1911091123"
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
            path: '/about/03_quickstart',
            meta: {
                h1: 'Quick Start',
                title: 'Quick Start',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-03_quickstart', '/about/03_quickstart', {})
        },
        {
            path: '/about/04_obtain_coherence_images',
            meta: {
                h1: 'Obtain Coherence Images',
                title: 'Obtain Coherence Images',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-04_obtain_coherence_images', '/about/04_obtain_coherence_images', {})
        },
        {
            path: '/about/05_cluster_discovery',
            meta: {
                h1: 'Coherence Cluster Discovery',
                title: 'Coherence Cluster Discovery',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('about-05_cluster_discovery', '/about/05_cluster_discovery', {})
        },
        {
            path: '/install/01_installation',
            meta: {
                h1: 'Coherence Operator Installation',
                title: 'Coherence Operator Installation',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('install-01_installation', '/install/01_installation', {})
        },
        {
            path: '/install/02_pre_release_versions',
            meta: {
                h1: 'Accessing Pre-Release Versions',
                title: 'Accessing Pre-Release Versions',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('install-02_pre_release_versions', '/install/02_pre_release_versions', {})
        },
        {
            path: '/app-deployment/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: 'Application Deployment',
                keywords: 'oracle coherence, kubernetes, operator, Application Deployment',
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployment-010_overview', '/app-deployment/010_overview', {})
        },
        {
            path: '/app-deployment/020_packaging',
            meta: {
                h1: 'Packaging Applications',
                title: 'Packaging Applications',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployment-020_packaging', '/app-deployment/020_packaging', {})
        },
        {
            path: '/app-deployment/030_roles',
            meta: {
                h1: 'Using Application Roles',
                title: 'Using Application Roles',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployment-030_roles', '/app-deployment/030_roles', {})
        },
        {
            path: '/app-deployment/060_persistence',
            meta: {
                h1: 'Persistence',
                title: 'Persistence',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployment-060_persistence', '/app-deployment/060_persistence', {})
        },
        {
            path: '/app-deployment/090_rolling',
            meta: {
                h1: 'Rolling Upgrades',
                title: 'Rolling Upgrades',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('app-deployment-090_rolling', '/app-deployment/090_rolling', {})
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
            path: '/metrics/040_scraping',
            meta: {
                h1: 'Using Your Own Prometheus',
                title: 'Using Your Own Prometheus',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-040_scraping', '/metrics/040_scraping', {})
        },
        {
            path: '/metrics/050_dashboards',
            meta: {
                h1: 'Grafana Dashboards',
                title: 'Grafana Dashboards',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('metrics-050_dashboards', '/metrics/050_dashboards', {})
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
            path: '/logging/030_own',
            meta: {
                h1: 'Using Your Own Elasticsearch',
                title: 'Using Your Own Elasticsearch',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-030_own', '/logging/030_own', {})
        },
        {
            path: '/logging/040_dashboards',
            meta: {
                h1: 'Kibana Dashboards',
                title: 'Kibana Dashboards',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('logging-040_dashboards', '/logging/040_dashboards', {})
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
            path: '/management/020_management_over_rest',
            meta: {
                h1: 'Management over ReST',
                title: 'Management over ReST',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-020_management_over_rest', '/management/020_management_over_rest', {})
        },
        {
            path: '/management/030_heapdump',
            meta: {
                h1: 'Generating Heap Dumps',
                title: 'Generating Heap Dumps',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-030_heapdump', '/management/030_heapdump', {})
        },
        {
            path: '/management/040_visualvm',
            meta: {
                h1: 'Using VisualVM',
                title: 'Using VisualVM',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-040_visualvm', '/management/040_visualvm', {})
        },
        {
            path: '/management/050_console',
            meta: {
                h1: 'Accessing the Console',
                title: 'Accessing the Console',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-050_console', '/management/050_console', {})
        },
        {
            path: '/management/060_cohql',
            meta: {
                h1: 'Accessing CohQL',
                title: 'Accessing CohQL',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('management-060_cohql', '/management/060_cohql', {})
        },
        {
            path: '/clusters/010_introduction',
            meta: {
                h1: 'CoherenceCluster CRD Overview',
                title: 'CoherenceCluster CRD Overview',
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
            path: '/clusters/050_coherence',
            meta: {
                h1: 'Configure Coherence',
                title: 'Configure Coherence',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-050_coherence', '/clusters/050_coherence', {})
        },
        {
            path: '/clusters/052_coherence_config_files',
            meta: {
                h1: 'Coherence Config Files',
                title: 'Coherence Config Files',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-052_coherence_config_files', '/clusters/052_coherence_config_files', {})
        },
        {
            path: '/clusters/054_coherence_storage_enabled',
            meta: {
                h1: 'Storage Enabled or Disabled Roles',
                title: 'Storage Enabled or Disabled Roles',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-054_coherence_storage_enabled', '/clusters/054_coherence_storage_enabled', {})
        },
        {
            path: '/clusters/056_coherence_image',
            meta: {
                h1: 'Setting the Coherence Image',
                title: 'Setting the Coherence Image',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-056_coherence_image', '/clusters/056_coherence_image', {})
        },
        {
            path: '/clusters/058_coherence_management',
            meta: {
                h1: 'Coherence Management over ReST',
                title: 'Coherence Management over ReST',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-058_coherence_management', '/clusters/058_coherence_management', {})
        },
        {
            path: '/clusters/060_coherence_metrics',
            meta: {
                h1: 'Coherence Metrics',
                title: 'Coherence Metrics',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-060_coherence_metrics', '/clusters/060_coherence_metrics', {})
        },
        {
            path: '/clusters/062_coherence_persistence',
            meta: {
                h1: 'Coherence Persistence',
                title: 'Coherence Persistence',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-062_coherence_persistence', '/clusters/062_coherence_persistence', {})
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
            path: '/clusters/080_jvm',
            meta: {
                h1: 'Configure the JVM',
                title: 'Configure the JVM',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-080_jvm', '/clusters/080_jvm', {})
        },
        {
            path: '/clusters/085_safe_scaling',
            meta: {
                h1: 'Configure Safe Scaling',
                title: 'Configure Safe Scaling',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-085_safe_scaling', '/clusters/085_safe_scaling', {})
        },
        {
            path: '/clusters/090_ports_and_services',
            meta: {
                h1: 'Expose Ports and Services',
                title: 'Expose Ports and Services',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-090_ports_and_services', '/clusters/090_ports_and_services', {})
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
            path: '/clusters/110_volumes',
            meta: {
                h1: 'Configure Additional Volumes',
                title: 'Configure Additional Volumes',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-110_volumes', '/clusters/110_volumes', {})
        },
        {
            path: '/clusters/115_environment_variables',
            meta: {
                h1: 'Environment Variables',
                title: 'Environment Variables',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-115_environment_variables', '/clusters/115_environment_variables', {})
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
            path: '/clusters/130_pod_scheduling',
            meta: {
                h1: 'Configure Pod Scheduling',
                title: 'Configure Pod Scheduling',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-130_pod_scheduling', '/clusters/130_pod_scheduling', {})
        },
        {
            path: '/clusters/140_resource_constraints',
            meta: {
                h1: 'Container Resource Limits',
                title: 'Container Resource Limits',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-140_resource_constraints', '/clusters/140_resource_constraints', {})
        },
        {
            path: '/clusters/150_readiness_liveness',
            meta: {
                h1: 'Readiness & Liveness Probes',
                title: 'Readiness & Liveness Probes',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: true
            },
            component: loadPage('clusters-150_readiness_liveness', '/clusters/150_readiness_liveness', {})
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
            path: '/developer/09_useful',
            meta: {
                h1: 'Useful Info',
                title: 'Useful Info',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
            },
            component: loadPage('developer-09_useful', '/developer/09_useful', {})
        },
        {
            path: '/developer/01_introduction',
            meta: {
                h1: 'Coherence Operator Development',
                title: 'Coherence Operator Development',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
            },
            component: loadPage('developer-01_introduction', '/developer/01_introduction', {})
        },
        {
            path: '/developer/04_how_it_works',
            meta: {
                h1: 'How It Works',
                title: 'How It Works',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
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
                hasNav: false
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
                hasNav: false
            },
            component: loadPage('developer-06_debugging', '/developer/06_debugging', {})
        },
        {
            path: '/developer/03_high_level',
            meta: {
                h1: 'High Level Design',
                title: 'High Level Design',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
            },
            component: loadPage('developer-03_high_level', '/developer/03_high_level', {})
        },
        {
            path: '/developer/07_execution',
            meta: {
                h1: 'Execution',
                title: 'Execution',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
            },
            component: loadPage('developer-07_execution', '/developer/07_execution', {})
        },
        {
            path: '/examples/010_overview',
            meta: {
                h1: 'Overview',
                title: 'Overview',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
            },
            component: loadPage('examples-010_overview', '/examples/010_overview', {})
        },
        {
            path: '/developer/08_docs',
            meta: {
                h1: 'Building the Docs',
                title: 'Building the Docs',
                description: null,
                keywords: null,
                customLayout: null,
                hasNav: false
            },
            component: loadPage('developer-08_docs', '/developer/08_docs', {})
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
                { href: '/about/03_quickstart', title: 'Quick Start' },
                { href: '/about/04_obtain_coherence_images', title: 'Obtain Coherence Images' },
                { href: '/about/05_cluster_discovery', title: 'Coherence Cluster Discovery' }
            ]
        },
        {
            title: 'Installation',
            action: 'settings',
            group: '/install',
            items: [
                { href: '/install/01_installation', title: 'Coherence Operator Installation' },
                { href: '/install/02_pre_release_versions', title: 'Accessing Pre-Release Versions' }
            ]
        },
        {
            title: 'Application Deployment',
            action: 'extension',
            group: '/app-deployment',
            items: [
                { href: '/app-deployment/010_overview', title: 'Overview' },
                { href: '/app-deployment/020_packaging', title: 'Packaging Applications' },
                { href: '/app-deployment/030_roles', title: 'Using Application Roles' },
                { href: '/app-deployment/060_persistence', title: 'Persistence' },
                { href: '/app-deployment/090_rolling', title: 'Rolling Upgrades' }
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
                { href: '/metrics/040_scraping', title: 'Using Your Own Prometheus' },
                { href: '/metrics/050_dashboards', title: 'Grafana Dashboards' }
            ]
        },
        {
            title: 'Logging',
            action: 'donut_large',
            group: '/logging',
            items: [
                { href: '/logging/010_overview', title: 'Overview' },
                { href: '/logging/020_logging', title: 'Enabling Log Capture' },
                { href: '/logging/030_own', title: 'Using Your Own Elasticsearch' },
                { href: '/logging/040_dashboards', title: 'Kibana Dashboards' }
            ]
        },
        {
            title: 'Management and Diagnostics',
            action: 'favorite_outline',
            group: '/management',
            items: [
                { href: '/management/010_overview', title: 'Overview' },
                { href: '/management/020_management_over_rest', title: 'Management over ReST' },
                { href: '/management/030_heapdump', title: 'Generating Heap Dumps' },
                { href: '/management/040_visualvm', title: 'Using VisualVM' },
                { href: '/management/050_console', title: 'Accessing the Console' },
                { href: '/management/060_cohql', title: 'Accessing CohQL' }
            ]
        },
        {
            title: 'Coherence CRD Reference',
            action: 'widgets',
            group: '/clusters',
            items: [
                { href: '/clusters/010_introduction', title: 'CoherenceCluster CRD Overview' },
                { href: '/clusters/020_k8s_resources', title: 'Coherence Clusters K8s Resources' },
                { href: '/clusters/030_roles', title: 'Define Coherence Roles' },
                { href: '/clusters/040_replicas', title: 'Role Replica Count' },
                { href: '/clusters/050_coherence', title: 'Configure Coherence' },
                { href: '/clusters/052_coherence_config_files', title: 'Coherence Config Files' },
                { href: '/clusters/054_coherence_storage_enabled', title: 'Storage Enabled or Disabled Roles' },
                { href: '/clusters/056_coherence_image', title: 'Setting the Coherence Image' },
                { href: '/clusters/058_coherence_management', title: 'Coherence Management over ReST' },
                { href: '/clusters/060_coherence_metrics', title: 'Coherence Metrics' },
                { href: '/clusters/062_coherence_persistence', title: 'Coherence Persistence' },
                { href: '/clusters/070_applications', title: 'Configure Applications' },
                { href: '/clusters/080_jvm', title: 'Configure the JVM' },
                { href: '/clusters/085_safe_scaling', title: 'Configure Safe Scaling' },
                { href: '/clusters/090_ports_and_services', title: 'Expose Ports and Services' },
                { href: '/clusters/100_logging', title: 'Logging Configuration' },
                { href: '/clusters/110_volumes', title: 'Configure Additional Volumes' },
                { href: '/clusters/115_environment_variables', title: 'Environment Variables' },
                { href: '/clusters/120_annotations', title: 'Configure Pod Annotations' },
                { href: '/clusters/125_labels', title: 'Configure Pod Labels' },
                { href: '/clusters/130_pod_scheduling', title: 'Configure Pod Scheduling' },
                { href: '/clusters/140_resource_constraints', title: 'Container Resource Limits' },
                { href: '/clusters/150_readiness_liveness', title: 'Readiness & Liveness Probes' },
                { href: '/clusters/190_service_account', title: 'Kubernetes Service Account' },
                { href: '/clusters/200_private_repos', title: 'Using Private Image Registries' }
            ]
        },
        { divider: true },
        { header: 'Additional resources' },
        {
            title: 'Slack',
            action: 'fa-slack',
            href: 'https://join.slack.com/t/oraclecoherence/shared_invite/enQtNzcxNTQwMTAzNjE4LTJkZWI5ZDkzNGEzOTllZDgwZDU3NGM2YjY5YWYwMzM3ODdkNTU2NmNmNDFhOWIxMDZlNjg2MzE3NmMxZWMxMWE',
            target: '_blank'
        },
        {
            title: 'Coherence Community',
            action: 'people',
            href: 'https://coherence.java.net',
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