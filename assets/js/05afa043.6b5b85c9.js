(window.webpackJsonp=window.webpackJsonp||[]).push([[3],{69:function(e,r,t){"use strict";t.r(r),t.d(r,"frontMatter",(function(){return i})),t.d(r,"metadata",(function(){return a})),t.d(r,"toc",(function(){return c})),t.d(r,"default",(function(){return l}));var n=t(3),o=(t(0),t(93));const i={sidebar_position:3},a={unversionedId:"queries/sorting",id:"queries/sorting",isDocsHomePage:!1,title:"Sorting",description:"In this section we'll learn how to sort results",source:"@site/docs/queries/sorting.md",sourceDirName:"queries",slug:"/queries/sorting",permalink:"/dqlx/docs/queries/sorting",editUrl:"https://github.com/fenos/dqlx-docs/edit/master/docs/queries/sorting.md",version:"current",sidebarPosition:3,frontMatter:{sidebar_position:3},sidebar:"tutorialSidebar",previous:{title:"Filters",permalink:"/dqlx/docs/queries/filters"},next:{title:"Pagination",permalink:"/dqlx/docs/queries/pagination"}},c=[{value:"OrderAsc",id:"orderasc",children:[]},{value:"OrderDesc",id:"orderdesc",children:[]},{value:"Multiple Sorting",id:"multiple-sorting",children:[]}],s={toc:c};function l({components:e,...r}){return Object(o.b)("wrapper",Object(n.a)({},s,r,{components:e,mdxType:"MDXLayout"}),Object(o.b)("p",null,"In this section we'll learn how to sort results"),Object(o.b)("h3",{id:"orderasc"},"OrderAsc"),Object(o.b)("p",null,"To sort result in ascending order use the ",Object(o.b)("inlineCode",{parentName:"p"},"OrderAsc")," function"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-go"},'db.Query(dqlx.HasFn("name")).\n    OrderAsc("name")\n')),Object(o.b)("h3",{id:"orderdesc"},"OrderDesc"),Object(o.b)("p",null,"To sort result in descending order use the ",Object(o.b)("inlineCode",{parentName:"p"},"OrderDesc")," function"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-go"},'db.Query(dqlx.HasFn("name")).\n    OrderDesc("name")\n')),Object(o.b)("h3",{id:"multiple-sorting"},"Multiple Sorting"),Object(o.b)("p",null,"You can chain multiple sorting criteria"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-go"},'db.Query(dqlx.HasFn("name")).\n    OrderDesc("name").\n    OrderAsc("age")\n')))}l.isMDXComponent=!0},93:function(e,r,t){"use strict";t.d(r,"a",(function(){return d})),t.d(r,"b",(function(){return m}));var n=t(0),o=t.n(n);function i(e,r,t){return r in e?Object.defineProperty(e,r,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[r]=t,e}function a(e,r){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);r&&(n=n.filter((function(r){return Object.getOwnPropertyDescriptor(e,r).enumerable}))),t.push.apply(t,n)}return t}function c(e){for(var r=1;r<arguments.length;r++){var t=null!=arguments[r]?arguments[r]:{};r%2?a(Object(t),!0).forEach((function(r){i(e,r,t[r])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):a(Object(t)).forEach((function(r){Object.defineProperty(e,r,Object.getOwnPropertyDescriptor(t,r))}))}return e}function s(e,r){if(null==e)return{};var t,n,o=function(e,r){if(null==e)return{};var t,n,o={},i=Object.keys(e);for(n=0;n<i.length;n++)t=i[n],r.indexOf(t)>=0||(o[t]=e[t]);return o}(e,r);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)t=i[n],r.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var l=o.a.createContext({}),u=function(e){var r=o.a.useContext(l),t=r;return e&&(t="function"==typeof e?e(r):c(c({},r),e)),t},d=function(e){var r=u(e.components);return o.a.createElement(l.Provider,{value:r},e.children)},p={inlineCode:"code",wrapper:function(e){var r=e.children;return o.a.createElement(o.a.Fragment,{},r)}},b=o.a.forwardRef((function(e,r){var t=e.components,n=e.mdxType,i=e.originalType,a=e.parentName,l=s(e,["components","mdxType","originalType","parentName"]),d=u(t),b=n,m=d["".concat(a,".").concat(b)]||d[b]||p[b]||i;return t?o.a.createElement(m,c(c({ref:r},l),{},{components:t})):o.a.createElement(m,c({ref:r},l))}));function m(e,r){var t=arguments,n=r&&r.mdxType;if("string"==typeof e||n){var i=t.length,a=new Array(i);a[0]=b;var c={};for(var s in r)hasOwnProperty.call(r,s)&&(c[s]=r[s]);c.originalType=e,c.mdxType="string"==typeof e?e:n,a[1]=c;for(var l=2;l<i;l++)a[l]=t[l];return o.a.createElement.apply(null,a)}return o.a.createElement.apply(null,t)}b.displayName="MDXCreateElement"}}]);