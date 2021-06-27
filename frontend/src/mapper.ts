import {
  PolicyReport, ClusterPolicyReport, GlobalPolicyReportMap, NamespacePolicyMap, Priority, Result,
} from '@/models';

const priorityToColor: { [key in Priority]: string } = {
  [Priority.SUCCESS]: 'green lighten-2',
  [Priority.DEBUG]: 'light-blue lighten-2',
  [Priority.INFO]: 'green darken-1',
  [Priority.WARNING]: 'orange lighten-1',
  [Priority.ERROR]: 'red darken-3',
  [Priority.CRITICAL]: 'red darken-4',
};

export const mapPriority = (priority: Priority): string => priorityToColor[priority] || priorityToColor[Priority.DEBUG];

export const convertPolicyReports = (reports: PolicyReport[]): NamespacePolicyMap => {
  const unordnered = reports.reduce<NamespacePolicyMap>((acc, report) => {
    if (acc.hasOwnProperty(report.namespace) === false) {
      acc[report.namespace] = {
        namespace: report.namespace,
        summary: {
          skip: 0, pass: 0, warn: 0, fail: 0, error: 0,
        },
        results: [],
      };
    }

    report.results.forEach((result) => {
      acc[report.namespace].summary[result.status] += 1;
      acc[report.namespace].results.push(result);
    });

    return acc;
  }, {});

  return Object.keys(unordnered).sort().reduce<NamespacePolicyMap>((acc, key) => {
    acc[key] = unordnered[key];

    return acc;
  }, {});
};

export const generateGlobalPolicyReports = (reports: PolicyReport[]): GlobalPolicyReportMap => {
  const unordnered = reports.reduce<GlobalPolicyReportMap>((acc, report) => {
    report.results.forEach((result) => {
      if (acc.hasOwnProperty(result.policy) === false) {
        acc[result.policy] = {
          summary: {
            skip: 0, pass: 0, warn: 0, fail: 0, error: 0,
          },
          results: [],
        };
      }

      acc[result.policy].summary[result.status] += 1;
      acc[result.policy].results.push(result);
    });

    return acc;
  }, {});

  return Object.keys(unordnered).sort().reduce<GlobalPolicyReportMap>((acc, key) => {
    acc[key] = unordnered[key];

    return acc;
  }, {});
};

export const findSources = (reports: PolicyReport[]): string[] => {
  const unordnered = reports.reduce<{ [source: string]: boolean }>((acc, report) => {
    report.results.forEach((result) => {
      if (result.source && acc.hasOwnProperty(result.source) === false) {
        acc[result.source] = true;
      }
    });

    return acc;
  }, {});

  return Object.keys(unordnered).sort();
};

export const flatPolicies = (reports: Array<PolicyReport|ClusterPolicyReport>) => reports.reduce<string[]>((acc, item) => {
  item.results.forEach((result: Result) => {
    if (acc.includes(result.policy)) {
      return;
    }

    acc.push(result.policy);
  });

  return acc;
}, []);
