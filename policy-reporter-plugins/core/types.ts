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
    id: string;
    kyverno: boolean;
}

export type Config = {
    plugins: string[];
    displayMode: DisplayMode | '';
    refreshInterval: number;
    views: ViewsCofig;
    clusters: Cluster[];
    labelFilter: string[];
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
    properties: {[key: string]: string};
}

export type Target = {
    name: string;
    minimumPriority: string;
    sources?: string[];
    skipExistingOnStartup: boolean;
}

export type StatusCount = {
    status: Status;
    count: number;
}

export type NamespacedStatusCount = {
    status: Status;
    items: Array<{ namespace: string; count: number; }>
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
    properties?: {[key: string]: string};
    scored: boolean;
    resource?: Resource;
}

export type ResultList = { items: ListResult[], count: number }

export interface CoreAPI {
    config(): Promise<Config>
    logs(): Promise<Result[]>
    targets(): Promise<Target[]>
    namespacedCategories(source?: string): Promise<string[]>
    namespaces(source?: string): Promise<string[]>
    namespacedPolicies(source?: string): Promise<string[]>
    namespacedRules(source?: string): Promise<string[]>
    namespacedSources(): Promise<string[]>
    namespacedKinds(source?: string): Promise<string[]>
    namespacedResults(filter?: Filter, pagination?: Pagination): Promise<ResultList>
    namespacedReportLabels(source?: string): Promise<{[key: string]: string[]}>
    namespacedStatusCount(filter?: Filter): Promise<NamespacedStatusCount[]>
    clusterPolicies(source?: string): Promise<string[]>
    clusterRules(source?: string): Promise<string[]>
    clusterKinds(source?: string): Promise<string[]>
    clusterSources(): Promise<string[]>
    clusterCategories(source?: string): Promise<string[]>
    clusterReportLabels(source?: string): Promise<{[key: string]: string[]}>
    statusCount(filter?: Filter): Promise<StatusCount[]>
    ruleStatusCount (policy: string, rule: string): Promise<StatusCount[]>
    results(filter?: Filter, pagination?: Pagination): Promise<ResultList>
    setPrefix (prefix: string): void
}

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
