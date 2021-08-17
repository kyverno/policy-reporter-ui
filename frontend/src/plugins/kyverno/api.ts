import subpath from '@/subpath';
import axios from 'axios';
import { Policy } from './models';

const production = process.env.NODE_ENV === 'production';

const instance = axios.create({ baseURL: production ? `//${window.location.host}${subpath()}api/kyverno` : `${process.env.VUE_APP_API}/kyverno` });

export default {
  async policies(): Promise<Policy[]> {
    const { data } = await instance.get('policies');

    return data;
  },
};
