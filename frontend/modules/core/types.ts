export enum Severity {
    INFO = 'info',
    LOW = 'low',
    MEDIUM = 'medium',
    HIGH = 'high',
    CRITICAL = 'critical',
}

export enum Priority {
    SUCCESS = 'success',
    DEBUG = 'debug',
    INFO = 'info',
    WARNING = 'warning',
    ERROR = 'error',
    CRITICAL = 'critical'
}

export enum Status {
    SKIP = 'skip',
    PASS = 'pass',
    WARN = 'warn',
    FAIL = 'fail',
    ERROR = 'error'
}

export enum DisplayMode {
    UNSPECIFIED = '',
    DARK = 'dark',
    LIGHT = 'light'
}

export enum Direction {
    ASC = 'asc',
    DESC = 'desc'
}

export type Dictionary = { [key: string]: string }

export type DashboardConfig = {
    policyReports: boolean;
    clusterPolicyReports: boolean;
}

export type ViewsCofig = {
    logs: boolean;
    policyReports: boolean;
    clusterPolicyReports: boolean;
    kyvernoPolicies: boolean;
    kyvernoVerifyImages: boolean;
    dashboard: DashboardConfig
}

export type Cluster = {
    name: string;
    slug: string;
    plugins: string[];
}

export type Category = {
    name: string;
    pass: number;
    warn: number;
    fail: number;
    error: number;
    skip: number;
}

export type PolicyResult = {
    name: string;
    namespace?: string;
    title: string;
    source: string;
    description: string;
    category: string;
    severity?: Severity;
    results: {
        pass?: number;
        warn?: number;
        fail?: number;
        error?: number;
        skip?: number;
    };
}

export type CustomBoard = {
    name: string;
    id: string;
}

export type CustomBoardDetails = {
    name: string;
    id: string;
    namespaces: string[];
    sources: string[];
    labels: { [label: string]: string };
}

export type Source = {
    name: string;
    categories: Category[];
}

export type SourceConfig = {
    name: string;
    dashboard: boolean;
    excludes: {
        namespaceKinds: string[]
        clusterKinds: string[]
    };
}

export type Profile = {
    id?: string;
    name: string;
}

export type Navigation = {
    title: string;
    subtitle: string;
    path: string;
    children?: Navigation[];
    exact?: boolean;
}

export type LayoutConfig = {
    targets: boolean;
    profile?: Profile;
    sources: Navigation[];
    policies: Navigation[];
    customBoards: Navigation[];
}

export type Dataset = {
    data: number[];
    backgroundColor: string | string[];
    label?: string;
}

export type Chart = {
    labels: string[];
    datasets: Dataset[];
    name: string;
}

export type ClusterScope = {
    [source: string]: {
        [key in Status]: number;
    }
}

export type NamespaceScope = {
    [source: string]: Chart
}

export type Findings = Chart

export type Dashboard = {
    title?: string;
    clusterScope: boolean;
    filterSources: string[];
    multiSource: boolean;
    showResults: string[];
    singleSource: boolean;
    charts: {
        clusterScope: ClusterScope;
        namespaceScope: {
            preview?: NamespaceScope
            complete: NamespaceScope
        };
        findings: { [key in Status]: Findings } | Findings
    };
    namespaces: string[];
    sources: string[];
    sourcesNavi: Array<{ title: string; name: string }>
    total: {
        count: number;
        perResult: {
            [key in Partial<Status>]: number;
        }
    }
}

export type SourceDetails = {
    title: string;
    name: string;
    categories: string[];
    chart: Chart;
}
export type ResourceDetails = {
    resource: Resource;
    results: { [key in Status]: number; }
    chart?: Chart
    sources: SourceDetails[]
}

export type PolicyDetails = {
    title: string;
    name: string;
    namespaces: string[];
    references?: string[];
    description: string;
    showDetails: boolean;
    engine?: {
        name: string;
        kubernetesVersion: string;
        version: string;
        subjects: string[];
    };
    sourceCode?: {
        contentType: string;
        content: string;
    };
    charts: {
        findings: Chart;
        namespaceScope: {
            complete: Chart;
            preview: Chart;
        };
        clusterScope: { [key in Status]: number; };
    };
    details: { title: string; value: string }[]
    additional: { title: string; items: { title?: string; value: string }[] }[];
}

export type Config = {
    error?: Error;
    plugins: string[];
    displayMode: DisplayMode;
    clusters: Cluster[];
    sources: SourceConfig[];
    oauth: boolean;
    default: string;
    defaultFilter: {
        resources: string[];
        clusterResources: string[];
    }
}

export type SourceFindings = { source: string; counts: { [key in Partial<Status>]: number }, total: number; }

export type FindingCounts = {
    total: number;
    counts: SourceFindings[]
}

export type ListResult = {
    id: string;
    namespace: string;
    kind: string;
    name: string;
    message: string;
    policy: string;
    rule: string;
    status: Status;
    properties: { [key: string]: string };
}

export type Target = {
    name: string;
    host?: string;
    skipTLS?: boolean;
    useTLS?: boolean;
    auth?: boolean;
    mountedSecret?: boolean;
    secretRef?: boolean;
    minimumPriority: Priority;
    filter: {
        namespaces?: { include: string[]; exclude: string[]; }
        priorities?: { include: string[]; exclude: string[]; }
        reportLabels?: { include: string[]; exclude: string[]; }
        policies?: { include: string[]; exclude: string[]; }
        sources?: { include: string[]; exclude: string[]; }
    };
    customFields: { [key: string]: string; }
    properties: { [key: string]: any; }
}

export type StatusCount = {
    status: Status;
    count: number;
}

export type NamespaceCounter = { namespaces: string[]; counts: number[] };

export type NamespacedStatusCount = {
    status: Status;
    items: Array<{ namespace: string; count: number; }>
}

export type NamespaceStatusCount = {
    [key: string]: { [status in Status]: number }
}

export type ResourceStatusCount = {
    [status in Status]: number;
} & {
    source: string;
}

export type Filter = {
    kinds?: string[];
    clusterKinds?: string[];
    categories?: string[];
    namespaces?: string[];
    severities?: Severity[];
    policies?: string[];
    sources?: string[];
    exclude?: string[];
    labels?: string[];
    status?: Status[];
    search?: string;
    resource_id?: string;
    namespaced?: boolean;
}

export type Pagination = {
    page: number;
    offset: number;
    sortBy?: string[]
    direction?: Direction
}

export type Resource = {
    apiVersion: string;
    kind: string;
    name: string;
    namespace?: string;
    uid: string;
    id: string;
}

export type Result = {
    message: string;
    policy: string;
    rule: string;
    priority: Priority;
    status: Status;
    source?: string;
    severity?: string;
    category?: string;
    properties?: { [key: string]: string };
    scored: boolean;
    resource?: Resource;
}

export type ResultList = { items: ListResult[], count: number }

export type ResourceResult = {
    id: string;
    uid: string;
    namespace: string;
    apiVersion: string;
    source: string;
    kind: string;
    name: string;
    pass: number;
    warn: number;
    fail: number;
    error: number;
    skip: number;
}

export type ResourceResultList = { items: ResourceResult[], count: number }

export type MappedResult = ListResult & { chips: Dictionary, cards: Dictionary, hasProps: boolean }

export type NamespaceCounters = { [status in Status]: { namespaces: string[]; counts: number[] } }
export type Counters = { [status in Status]: number }

export const createNamespaceCounters = (): NamespaceCounters => ({
    [Status.SKIP]: { namespaces: [], counts: [] },
    [Status.PASS]: { namespaces: [], counts: [] },
    [Status.WARN]: { namespaces: [], counts: [] },
    [Status.FAIL]: { namespaces: [], counts: [] },
    [Status.ERROR]: { namespaces: [], counts: [] }
})

export const createCounters = (): Counters => ({
    [Status.SKIP]: 0,
    [Status.PASS]: 0,
    [Status.WARN]: 0,
    [Status.FAIL]: 0,
    [Status.ERROR]: 0
})

export const createStatusList = (): Status[] => [
    Status.FAIL,
    Status.PASS,
    Status.WARN,
    Status.ERROR,
    Status.SKIP
]

export type NamespacedFilterAPI = 'namespacedKinds' | 'namespacedPolicies' | 'namespacedRules' | 'namespacedSources';

export type ClusterFilterAPI = 'clusterKinds' | 'clusterPolicies' | 'clusterRules' | 'clusterSources' | 'namespaces' | 'categories'
