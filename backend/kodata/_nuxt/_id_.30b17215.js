import{S as e,a as v,m as S,B as G,b as ee,e as F}from"./mapper.f869c811.js";import{d as b,c as O,t as u,O as m,T as Q,U as X,Q as l,ag as te,A as a,B as L,y as h,R as H,V as z,x as P,v as y,z as n,ab as K,a9 as Y,a0 as W,_ as C,ah as se,r as M,N as Z,H as j,Z as V,a5 as ae,a6 as re,ai as oe,aj as le,ak as ne,P as I,al as ue,X as w,Y as T,am as ce,W as E,a1 as q,$ as ie,L as de,an as pe}from"./entry.a7ce8945.js";import{_ as _e}from"./Search.vue.f75a0caa.js";const me=b({__name:"StatusPerCategory",props:{source:{}},setup(_){const c=_,t=O(()=>{const r={};c.source.categories.forEach(o=>{r[o.name||"other"]={[e.PASS]:o.pass,[e.SKIP]:o.skip,[e.FAIL]:o.fail,[e.WARN]:o.warn,[e.ERROR]:o.error}});const s=Object.keys(r).sort((o,p)=>o.localeCompare(p)).reduce((o,p)=>({...o,[p]:r[p]}),{}),d=Object.keys(s),i={[e.PASS]:{data:[],label:v(e.PASS),backgroundColor:S(e.PASS)},[e.FAIL]:{data:[],label:v(e.FAIL),backgroundColor:S(e.FAIL)},[e.WARN]:{data:[],label:v(e.WARN),backgroundColor:S(e.WARN)},[e.ERROR]:{data:[],label:v(e.ERROR),backgroundColor:S(e.ERROR)}};return d.forEach(o=>{i[e.PASS].data.push(s[o][e.PASS]),i[e.FAIL].data.push(s[o][e.FAIL]),i[e.WARN].data.push(s[o][e.WARN]),i[e.ERROR].data.push(s[o][e.ERROR])}),{style:{minHeight:`${125+d.length*25}px`},data:{labels:d,datasets:Object.values(i)},options:{height:"100%",indexAxis:"y",responsive:!0,maintainAspectRatio:!1,plugins:{title:{display:!0,text:"Results per Category"},legend:{display:!0,position:"bottom"}},scales:{x:{stacked:!0},y:{stacked:!0}}}}});return(r,s)=>(u(),m(l(G),Q(X(l(t))),null,16))}}),he=b({__name:"Status",props:{status:{type:String,default:null}},setup(_){const c=_,t=te(),r=O(()=>t.current.value.dark?ee(c.status):S(c.status));return(s,d)=>(u(),m(z,H({color:l(r),dark:""},s.$attrs),{default:a(()=>[L(h(_.status),1)]),_:1},16,["color"]))}}),fe=b({__name:"Severity",props:{severity:{type:String,default:null}},setup(_){const c=_,t=O(()=>{switch(c.severity){case F.INFO:return"info lighten-1";case F.LOW:return"info";case F.MEDIUM:return"warning";case F.HIGH:return"error";case F.CRITICAL:return"error darken-2";default:return"grey"}});return(r,s)=>(u(),m(z,H({color:l(t),dark:""},r.$attrs),{default:a(()=>[L(h(_.severity),1)]),_:1},16,["color"]))}}),ye={class:"inline-block mr-1"},ge=["href"],Re={key:1},ve=b({__name:"PropertyChip",props:{label:{},value:{}},setup(_){const c=t=>{let r;try{r=new URL(t)}catch{return!1}return r.protocol==="http:"||r.protocol==="https:"};return(t,r)=>(u(),m(z,{variant:"flat",pill:"",color:"indigo-lighten-1"},{default:a(()=>[P("span",ye,h(l(v)(t.label))+":",1),L(),c(t.value)?(u(),y("a",{key:0,href:t.value,target:"_blank",class:"text-white"},h(t.value),9,ge)):(u(),y("span",Re,h(t.value),1))]),_:1}))}}),ke={class:"text-subtitle-2"},Se=["href"],be=["innerHTML"],Ae=b({__name:"PropertyCard",props:{label:{},value:{}},setup(_){const c=t=>{let r;try{r=new URL(t)}catch{return!1}return r.protocol==="http:"||r.protocol==="https:"};return(t,r)=>(u(),m(C,H({variant:"outlined"},t.$attrs,{class:"bg-white"}),{default:a(()=>[n(K,{class:"bg-indigo"},{default:a(()=>[P("div",ke,h(t.label),1)]),_:1}),n(Y),n(W,{class:"text-body-2 pa-2"},{default:a(()=>[c(t.value)?(u(),y("a",{key:0,href:t.value,target:"_blank"},h(t.value),9,Se)):(u(),y("span",{key:1,innerHTML:t.value},null,8,be))]),_:1})]),_:1},16))}}),Pe={key:0},Ce={key:1},$e={class:"table-expand-text"},xe=["colspan"],Le={key:0},Oe={key:1},Ie=b({__name:"ResultTable",props:{source:{},category:{},resource:{},Status:{}},setup(_){const c=_,t=se({itemsPerPage:10,page:1,sortDesc:[],sortBy:[],groupBy:[],groupDesc:[],multiSort:!1,mustSort:!1}),r=M(!0),s=M(""),d=g=>Object.keys(g).sort().reduce((R,$)=>(R[$]=g[$],R),{}),{data:i,refresh:o}=Z(g=>g.results(c.resource,{page:t.page,offset:t.itemsPerPage},{sources:c.source?[c.source]:void 0,categories:c.category?[c.category]:void 0,search:s.value}),{default:()=>({items:[],count:0})});j(()=>t.page,o),j(()=>t.itemsPerPage,o),j(s,o);const p=O(()=>i.value.items.map(({properties:g,...R})=>{const $={},B={};let U=!1;for(const x in g)g[x].length>75?B[x]=g[x]:$[x]=g[x],U=!0;return{...R,properties:{},cards:d(B),chips:d($),hasProps:U}})),k=[{title:"Policy",key:"policy",width:"33%"},{title:"Rule",key:"rule",width:"33%"},{title:"Severity",key:"severity",width:"17%"},{title:"Status",key:"status",width:"17%"}];return(g,R)=>{const $=_e,B=he,U=fe,x=ve,J=Ae;return l(i).count>0||l(s)?(u(),m(E,{key:0},{default:a(()=>[n(V,{cols:"12"},{default:a(()=>[n(C,null,{title:a(()=>[g.category?(u(),y("span",Pe,h(g.category),1)):(u(),y("span",Ce,"Results for "+h(l(v)(g.source)),1))]),append:a(()=>[n($,{modelValue:l(s),"onUpdate:modelValue":R[0]||(R[0]=f=>ae(s)?s.value=f:null),class:"mr-4",style:{"min-width":"400px",float:"left",height:"48px"}},null,8,["modelValue"]),n(re,{icon:l(r)?"mdi-chevron-up":"mdi-chevron-down",onClick:R[1]||(R[1]=f=>r.value=!l(r)),variant:"text"},null,8,["icon"])]),default:a(()=>[n(oe,null,{default:a(()=>[le(P("div",null,[n(Y),n(ne,{items:l(p),"items-length":l(i).count,headers:k,"item-value":"id","show-expand":"","items-per-page":l(t).itemsPerPage,"onUpdate:itemsPerPage":R[2]||(R[2]=f=>l(t).itemsPerPage=f),page:l(t).page,"onUpdate:page":R[3]||(R[3]=f=>l(t).page=f)},{"item.status":a(({value:f})=>[n(B,{onClick:A=>s.value=f,status:f},null,8,["onClick","status"])]),"item.severity":a(({value:f})=>[f?(u(),m(U,{key:0,onClick:A=>s.value=f,severity:f},null,8,["onClick","severity"])):I("",!0)]),"expanded-row":a(({columns:f,item:A})=>[P("tr",$e,[P("td",{colspan:f.length,class:"py-3"},[A.hasProps?(u(),y("div",Le,[A.message?(u(),m(C,{key:0,variant:"flat"},{default:a(()=>[n(ue,{type:"info",variant:"flat"},{default:a(()=>[L(h(A.message),1)]),_:2},1024)]),_:2},1024)):I("",!0),P("div",null,[(u(!0),y(w,null,T(A.chips,(D,N)=>(u(),m(x,{key:N,label:N,value:D,class:"mr-2 mt-2 rounded-xl"},null,8,["label","value"]))),128)),(u(!0),y(w,null,T(A.cards,(D,N)=>(u(),m(J,{key:N,label:N,value:D,class:"mt-2"},null,8,["label","value"]))),128))])])):(u(),y("div",Oe,h(A.message),1))],8,xe)])]),_:1},8,["items","items-length","items-per-page","page"])],512),[[ce,l(r)]])]),_:1})]),_:1})]),_:1})]),_:1})):I("",!0)}}});const Ve={class:"text-h4"},we=b({__name:"CategoryTables",props:{source:{},resource:{}},setup(_){return(c,t)=>{const r=me,s=Ie;return u(),y(w,null,[n(E,null,{default:a(()=>[n(V,{cols:"12"},{default:a(()=>[n(C,null,{default:a(()=>[n(K,{class:"py-3"},{default:a(()=>[P("span",Ve,h(l(v)(c.source.name)),1)]),_:1})]),_:1})]),_:1})]),_:1}),n(E,null,{default:a(()=>[n(V,{cols:"12"},{default:a(()=>[n(C,null,{default:a(()=>[n(W,null,{default:a(()=>[n(r,{source:c.source},null,8,["source"])]),_:1})]),_:1})]),_:1})]),_:1}),(u(!0),y(w,null,T(c.source.categories,d=>(u(),m(s,{source:c.source.name,resource:c.resource.id,key:d.name,category:d.name},null,8,["source","resource","category"]))),128))],64)}}}),Ee=b({__name:"ResourceStatus",props:{data:{}},setup(_){const c=_,t=O(()=>{var o;const r={};(o=c.data)==null||o.forEach(p=>{p.items.forEach(k=>{r[k.source]||(r[k.source]={[e.PASS]:0,[e.SKIP]:0,[e.FAIL]:0,[e.WARN]:0,[e.ERROR]:0}),r[k.source][p.status]=k.count})});const s=Object.keys(r).sort((p,k)=>p.localeCompare(k)).reduce((p,k)=>({...p,[k]:r[k]}),{}),d=Object.keys(s),i={[e.PASS]:{data:[],label:v(e.PASS),backgroundColor:S(e.PASS)},[e.FAIL]:{data:[],label:v(e.FAIL),backgroundColor:S(e.FAIL)},[e.WARN]:{data:[],label:v(e.WARN),backgroundColor:S(e.WARN)},[e.ERROR]:{data:[],label:v(e.ERROR),backgroundColor:S(e.ERROR)}};return d.forEach(p=>{i[e.PASS].data.push(s[p][e.PASS]),i[e.FAIL].data.push(s[p][e.FAIL]),i[e.WARN].data.push(s[p][e.WARN]),i[e.ERROR].data.push(s[p][e.ERROR])}),{style:{minHeight:`${125+d.length*25}px`},data:{labels:d.map(p=>v(p)),datasets:Object.values(i)},options:{height:"100%",indexAxis:"y",responsive:!0,maintainAspectRatio:!1,plugins:{title:{display:!0,text:"Resource results per Source"},legend:{display:!0,position:"bottom"}},scales:{x:{stacked:!0},y:{stacked:!0}}}}});return(r,s)=>(u(),m(l(G),Q(X(l(t))),null,16))}}),Ne=b({__name:"ResourceResultCounts",props:{data:{}},setup(_){const c=_,t=O(()=>{var r;return(r=c.data)==null?void 0:r.reduce((s,d)=>(d.status===e.SKIP||(s[d.status]=d.items.reduce((i,o)=>i+o.count,0)),s),{[e.PASS]:0,[e.WARN]:0,[e.FAIL]:0,[e.ERROR]:0})});return(r,s)=>(u(),m(q,{fluid:""},{default:a(()=>[n(E,null,{default:a(()=>[(u(!0),y(w,null,T(l(t),(d,i)=>(u(),m(V,{key:i,cols:"12",sm:"6",md:"3"},{default:a(()=>[n(C,{flat:"",title:`${i} results`,class:"text-white text-center",style:ie(`background-color: ${l(S)(i)}`)},{default:a(()=>[n(W,{class:"text-h3 my-4"},{default:a(()=>[L(h(d),1)]),_:2},1024)]),_:2},1032,["title","style"])]),_:2},1024))),128))]),_:1})]),_:1}))}}),Fe={class:"bg-indigo"},We={key:0},De=b({__name:"[id]",setup(_){const c=de(),{data:t}=Z(s=>Promise.all([s.resource(c.params.id),s.resourceStatusCount(c.params.id),s.sources(c.params.id)]).then(([d,i,o])=>({resource:d,counts:i,sources:o})),{default:()=>({resource:{},counts:[],sources:[]})}),r=O(()=>{var s;return(((s=t.value)==null?void 0:s.sources)||[]).sort((d,i)=>d.name.localeCompare(i.name))});return(s,d)=>{const i=we;return l(t)?(u(),m(q,{fluid:"",class:"py-4 px-4 main-height",key:l(c).params.id},{default:a(()=>[n(E,null,{default:a(()=>[n(V,null,{default:a(()=>[n(C,{elevation:"2",rounded:""},{default:a(()=>[P("div",Fe,[n(K,null,{default:a(()=>{var o,p;return[(o=l(t))!=null&&o.resource.namespace?(u(),y("span",We,h((p=l(t))==null?void 0:p.resource.namespace)+"/",1)):I("",!0),L(h(l(t).resource.name),1)]}),_:1}),n(pe,{class:"pb-4"},{default:a(()=>{var o;return[L(h((o=l(t))==null?void 0:o.resource.apiVersion)+" "+h(l(t).resource.kind),1)]}),_:1})]),n(W,null,{default:a(()=>[n(Ne,{data:l(t).counts},null,8,["data"])]),_:1})]),_:1})]),_:1})]),_:1}),l(t).sources.length>1?(u(),m(E,{key:0},{default:a(()=>[n(V,null,{default:a(()=>[n(C,null,{default:a(()=>[n(W,null,{default:a(()=>[n(Ee,{data:l(t).counts},null,8,["data"])]),_:1})]),_:1})]),_:1})]),_:1})):I("",!0),(u(!0),y(w,null,T(l(r),o=>(u(),m(i,{source:o,resource:l(t).resource,key:o},null,8,["source","resource"]))),128))]),_:1})):I("",!0)}}});export{De as default};
