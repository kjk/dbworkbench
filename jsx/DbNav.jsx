import React from 'react';
import action from './action.js';
import view from './view.js';

export default class DbNav extends React.Component {
  render() {
    //console.log("DbNav.render: view: ", this.props.view);
    const currentView = this.props.view;
    const children = view.MainTabViews.map(function(viewName) {
      const handler = function() {
        action.viewSelected(viewName);
      };

      const selected = (currentView == viewName);
      if (selected) {
        return <li key={viewName} onClick={handler} className="selected"><u>{viewName}</u></li>;
      } else {
        return <li key={viewName} onClick={handler}>{viewName}</li>;
      }
    });

    return (
      <div id="nav">
        <ul>
          {children}
        </ul>
      </div>
    );
  }
}
