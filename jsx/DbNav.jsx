import React from 'react';
import * as action from './action.js';
import * as api from './api.js';
import view from './view.js';
import Modal from 'react-modal';
import { OverlayTrigger, Tooltip } from 'react-bootstrap';


export default class DbNav extends React.Component {
  constructor(props, context) {
    super(props, context);

    this.handleFeedbackButton = this.handleFeedbackButton.bind();
  }

  handleFeedbackButton(e) {
    e.preventDefault();
    api.launchBrowserWithURL("http://dbheroapp.com/feedback");
  }

  render() {
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

    const tooltip = (
      <Tooltip id="feedback">Let us know how we can improve dbHero.</Tooltip>
    );

    return (
      <div id="nav">
        <ul>
          {children}
        </ul>
        <OverlayTrigger placement="left" overlay={tooltip}>
          <button className="feedback-button" onClick={this.handleFeedbackButton.bind(this)}>Feedback</button>
        </OverlayTrigger>
      </div>
    );
  }
}

