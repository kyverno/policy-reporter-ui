export enum Severity {
    UNKNOWN = 'unknown',
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
    ERROR = 'error',
    SUMMARY = 'summary'
}

export enum DisplayMode {
    UNSPECIFIED = '',
    DARK = 'dark',
    LIGHT = 'light',
    COLOR_BLIND = 'colorblind',
    COLOR_BLIND_DARK = 'colorblinddark'
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
    type?: string;
}

export type ClusterScope = {
    [source: string]: {
        [key in Status | Severity]: number;
    }
}

export type NamespaceScope = {
    [source: string]: Chart
}

export type Findings = Chart

export type ViewType = 'severity' | 'status' | ''

export type Dashboard = {
    title?: string;
    clusterScope: boolean;
    filterSources: string[];
    multiSource: boolean;
    showResults: string[];
    status: Status[];
    severities: Severity[];
    singleSource: boolean;
    exceptions: boolean;
    type: ViewType;
    charts: {
        clusterScope: ClusterScope;
        namespaceScope: {
            [source: string]: {
                preview?: Chart
                complete: Chart
            }
        };
        findings: { [key in Status]: Findings } | Findings
    };
    namespaces: string[];
    sources: string[];
    sourcesNavi: Array<{ title: string; name: string }>
    total: {
        count: number;
        perResult: {
            [key in Partial<Status | Severity>]: number;
        }
    }
}

export type SourceDetails = {
    title: string;
    name: string;
    status: Status[];
    categories: string[];
    chart: Chart;
    exceptions: boolean;
    plugin: boolean;
}

export type ResourceDetails = {
    resource: Resource;
    results: { [key in Status]: number; }
    severityResults: { [key in Severity]: number; }
    chart?: Chart
    sources: SourceDetails[]
    status: Status[];
    severities: Severity[];
}

export type PolicyDetails = {
    title: string;
    name: string;
    namespaces: string[];
    references?: string[];
    description: string;
    showDetails: boolean;
    exceptions: boolean;
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
    resourceId: string;
    message: string;
    policy: string;
    rule: string;
    status: Status;
    severity: Severity;
    timestamp: number;
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
    status: {
        pass: number;
        warn: number;
        fail: number;
        error: number;
        skip: number;
    }
    severities: {
        unknown: number;
        info: number;
        low: number;
        medium: number;
        high: number;
        critical: number;
    }
}

export type ExceptionRule = { name: string, props: { [key: string]: string } }
export type ExceptionPolicy = { name: string, rules: ExceptionRule[] }

export type ExceptionResponse = {
    minVersion?: string;
    resource: string;
}

export type ResourceResultList = { items: ResourceResult[], count: number }

export type MappedResult = ListResult & { chips: Dictionary, cards: Dictionary, hasProps: boolean }

export type NamespaceCounters = Omit<{ [status in Status]: { namespaces: string[]; counts: number[] } }, Status.SUMMARY>
export type Counters = Omit<{ [status in Status]: number }, Status.SUMMARY>

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
