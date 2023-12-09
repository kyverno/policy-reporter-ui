import{d as I,r as _,t as l,v as y,z as t,a9 as L,aM as B,A as a,a6 as N,Q as s,aA as S,aI as g,B as R,y as f,aJ as P,aG as $,O as d,P as v,X as k,ae as D,af as E,aN as O,H as T,M as C,aD as z,Y as M,a0 as b,_ as q,q as F}from"./entry.a7ce8945.js";import{_ as j}from"./ResultChip.vue.527b52bf.js";import{S as m}from"./mapper.f869c811.js";const G=D(()=>E(()=>import("./ResourceResults.8b905684.js"),["./ResourceResults.8b905684.js","./ResultChip.vue.527b52bf.js","./mapper.f869c811.js","./entry.a7ce8945.js","./entry.729f9978.css"],import.meta.url).then(e=>e.default||e)),H=I({__name:"ResourceResultItem",props:{item:{type:Object,required:!0},details:{type:Boolean,default:!1}},setup(e){const n=_(!1);return(p,u)=>{const o=j,r=G;return l(),y(k,null,[t(L),t($,{to:`/resource/${e.item.id}`},B({append:a(()=>[t(o,{status:s(m).PASS,count:e.item.pass,tooltip:"pass results"},null,8,["status","count"]),t(o,{class:"ml-2",status:s(m).WARN,count:e.item.warn,tooltip:"warning results"},null,8,["status","count"]),t(o,{class:"ml-2",status:s(m).FAIL,count:e.item.fail,tooltip:"fail results"},null,8,["status","count"]),t(o,{class:"ml-2",status:s(m).ERROR,count:e.item.error,tooltip:"error results"},null,8,["status","count"])]),default:a(()=>[t(g,null,{default:a(()=>[R(f(e.item.name)+" ",1)]),_:1}),t(P,null,{default:a(()=>[R(f(e.item.apiVersion)+" "+f(e.item.kind),1)]),_:1})]),_:2},[e.details?{name:"prepend",fn:a(()=>[t(N,{class:"mr-2",variant:"text",icon:s(n)?"mdi-chevron-down":"mdi-chevron-up",onClick:u[0]||(u[0]=S(i=>n.value=!s(n),["stop","prevent"]))},null,8,["icon"])]),key:"0"}:void 0]),1032,["to"]),s(n)?(l(),d(r,{key:0,id:e.item.id},null,8,["id"])):v("",!0)],64)}}}),Y=I({__name:"ResourceResultList",props:{namespace:{},details:{type:Boolean}},async setup(e){let n,p;const u=e,{$coreAPI:o}=F(),r=_(!0),i=_({items:[],count:0}),V=async()=>{try{i.value=await o.namespacedResourceResults({namespaces:[u.namespace],kinds:C.value})}catch(c){console.error(c)}finally{r.value=!1}};return[n,p]=O(()=>V()),await n,p(),T(C,V),(c,J)=>{const x=H;return l(),d(q,{title:c.namespace},{default:a(()=>{var h,w;return[(w=(h=s(i))==null?void 0:h.items)!=null&&w.length?(l(),d(z,{key:0,lines:"two"},{default:a(()=>[(l(!0),y(k,null,M(s(i).items,A=>(l(),d(x,{key:A.id,item:A,details:c.details},null,8,["item","details"]))),128))]),_:1})):v("",!0),!s(r)&&!s(i).items.length?(l(),y(k,{key:1},[t(L),t(b,null,{default:a(()=>[R(" No resources for the selected kinds found ")]),_:1})],64)):v("",!0)]}),_:1},8,["title"])}}});export{Y as default};
