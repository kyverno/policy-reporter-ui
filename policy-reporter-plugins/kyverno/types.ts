import { Result, Severity } from '../core/types'

export enum RuleType {
    VALIDATION = 'validation',
    MUTATION = 'mutation',
    GENERATION = 'generation',
}

export enum ReportType {
    POLICY = 'Policy Report',
    NAMESPACE = 'Namespace Report'
}

export interface Rule {
    message?: string;
    name: string;
    type: RuleType;
}

export interface Policy {
    uid: string;
    kind: string;
    namespace?: string;
    category?: string;
    description?: string;
    severity?: Severity;
    autogenControllers: string[];
    name: string;
    rules: Rule[];
    validationFailureAction: string;
    background: boolean;
    content: string;
}

export interface VerifyImageRule {
    rule: string;
    attestations?: string;
    policy: { name: string; namespace?: string; uid: string }
    repository: string;
    image: string;
    key: string;
}

export type PolicyGroups = { [category: string]: Policy[] };

export type ResultMap = {
  pass: Result[];
  fail: Result[];
}

export type Filter = {
    policies: string[];
    namespaces: string[];
    clusterScope: 1 | 0
}

export interface KyvernoAPI {
    policies (): Promise<{ policies: Policy[]; groups: PolicyGroups }>;
    verifyImageRules (): Promise<VerifyImageRule[]>;
    policyReport (filter: Filter): Promise<BlobPart>;
    namespaceReport (filter: Filter): Promise<BlobPart>;
    setPrefix (prefix: string): void;
}
