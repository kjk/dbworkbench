import React from 'react';
import { Td } from './td.jsx';
import { toArray, filterPropsFrom } from './utils.jsx';

function toTdChildren(children, data, columns) {
  if (!data || !columns || typeof columns.map !== 'function') {
    return children;
  }
  children = children.concat(columns.map((column, i) => {
    if (!data.hasOwnProperty(column.key)) {
      return <Td column={ column } key={ column.key } />;
    }

    let value = data[column.key];
    let props = {};

    if (value && value.__reactableMeta === true) {
      props = value.props;
      value = value.value;
    }

    return <Td column={ column } key={ column.key } {...props}>
             { value }
           </Td>;
  }));
  return children;
}

export class Tr extends React.Component {
  render() {
    let children = React.Children.toArray();
    children = toTdChildren(children, this.props.data, this.props.columns);
    var props = filterPropsFrom(this.props);
    return React.DOM.tr(props, children);
  }
}

Tr.childNode = Td;
Tr.dataType = 'object';

