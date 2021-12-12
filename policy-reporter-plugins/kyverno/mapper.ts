import { Policy } from './types'

export default (policy: Partial<Policy>) => ({
  ...policy,
  autogenControllers: policy.autogenControllers || []
})
