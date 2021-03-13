import axios from 'axios';

import { PolicyReport, ClusterPolicyReport, Target } from '@/models';

const production = process.env.NODE_ENV === 'production';

const instance = axios.create({ baseURL: production ? `//${window.location.host}/api` : process.env.VUE_APP_API });

export default {
  async policyReports(): Promise<PolicyReport[]> {
    const { data } = await instance.get('policy-reports');

    return data;
  },
  async clusterPolicyReports(): Promise<ClusterPolicyReport[]> {
    const { data } = await instance.get('cluster-policy-reports');

    return data;
  },
  async targets(): Promise<Target[]> {
    const { data } = await instance.get('targets');

    return data;
  },
};
