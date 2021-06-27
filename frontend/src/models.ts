export type Summary = {
    pass: number;
    skip: number;
    warn: number;
    error: number;
    fail: number;
}

export type Resource = {
    apiVersion: string;
    kind: string;
    name: string;
    namespace?: string;
    uid: string;
}

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

export type PolicyReport = {
    name: string;
    namespace: string;
    summary: Summary;
    results: Result[];
    creationTimesamp: Date;
}

export type NamespacePolicyReport = {
    namespace: string;
    summary: Summary;
    results: Result[];
}

export type GlobalPolicyReport = {
    summary: Summary;
    results: Result[];
}

export type ClusterPolicyReport = {
    name: string;
    summary: Summary;
    results: Result[];
    creationTimesamp: Date;
}

export type Target = {
    name: string;
    minimumPriority: string;
    skipExistingOnStartup: boolean;
}

export type NamespacePolicyMap = { [namespace: string]: NamespacePolicyReport };
export type GlobalPolicyReportMap = { [policy: string]: GlobalPolicyReport };
