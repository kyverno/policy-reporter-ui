export enum Severity {
    LOW = 'low',
    MEDIUM = 'medium',
    HIGH = 'high',
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

export type Config = {
    plugins: string[];
    displayMode: DisplayMode | '';
    views: ViewsCofig;
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
    categories(source?: string): Promise<string[]>
    clusterKinds(source?: string): Promise<string[]>
    namespacedPolicies(source?: string): Promise<string[]>
    namespacedRules(source?: string): Promise<string[]>
    clusterPolicies(source?: string): Promise<string[]>
    clusterRules(source?: string): Promise<string[]>
    namespaces(source?: string): Promise<string[]>
    clusterSources(): Promise<string[]>
    namespacedSources(): Promise<string[]>
    namespacedKinds(source?: string): Promise<string[]>
    statusCount(filter?: Filter): Promise<StatusCount[]>
    ruleStatusCount (policy: string, rule: string): Promise<StatusCount[]>
    namespacedStatusCount(filter?: Filter): Promise<NamespacedStatusCount[]>
    results(filter?: Filter, pagination?: Pagination): Promise<ResultList>
    namespacedResults(filter?: Filter, pagination?: Pagination): Promise<ResultList>
}
