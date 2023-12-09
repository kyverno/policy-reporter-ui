export enum Severity {
    INFO = 'low',
    LOW = 'low',
    MEDIUM = 'medium',
    HIGH = 'high',
    CRITICAL = 'high',
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

export type Source = {
    name: string;
    categories: Category[];
}

export type Config = {
    plugins: string[];
    displayMode: DisplayMode | '';
    views: ViewsCofig;
    clusters: Cluster[];
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
    minimumPriority: Priority;
    sources?: string[];
    skipExistingOnStartup: boolean;
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

export type ResourceStatusCount = {
    status: Status;
    items: Array<{ source: string; count: number; }>
}

export type Filter = {
    kinds?: string[];
    categories?: string[];
    namespaces?: string[];
    severities?: Severity[];
    policies?: string[];
    sources?: string[];
    status?: Status[];
    search?: string;
    resource_id?: string;
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
