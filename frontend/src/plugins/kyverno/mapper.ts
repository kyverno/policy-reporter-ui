import { Policy } from './models';

export default (policy: Partial<Policy>) => ({
  ...policy,
  autogenControllers: policy.autogenControllers || [],
});
