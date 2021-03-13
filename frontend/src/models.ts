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

export enum Priority {
    DEBUG = 'debug',
    INFO = 'info',
    WARNING = 'warning',
    ERROR = 'error'
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
    severity?: string;
    category?: string;
    scored: boolean;
    resource: Resource;
}

export type PolicyReport = {
    name: string;
    namespace: string;
    summary: Summary;
    results: Result[];
    creationTimesamp: Date;
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

const priorityToColor: { [key in Priority]: string } = {
  [Priority.DEBUG]: 'light-blue lighten-2',
  [Priority.INFO]: 'green darken-1',
  [Priority.WARNING]: 'orange lighten-1',
  [Priority.ERROR]: 'red darken-3',
};

export const mapPriority = (priority: Priority): string => priorityToColor[priority] || priorityToColor[Priority.DEBUG];
