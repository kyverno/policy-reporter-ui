import { Result, Severity } from '../core/types'

export enum RuleType {
    VALIDATION = 'validation',
    MUTATION = 'mutation',
    GENERATION = 'generation',
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

export type PolicyGroups = { [category: string]: Policy[] };

export type ResultMap = {
  pass: Result[];
  fail: Result[];
}

export interface KyvernoAPI {
    policies (): Promise<{ policies: Policy[]; groups: PolicyGroups }>
}
