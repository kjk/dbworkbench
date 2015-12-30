import React from 'react';
import * as action from './action.js';
import view from './view.js';
import Modal from 'react-modal';

export default class DbNav extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.closeModal = this.closeModal.bind(this);

    this.state = {
      modalIsOpen: false,
    };
  }

  closeModal() {
    this.setState({modalIsOpen: false});
  }

  handleModalCloseRequest() {
    // opportunity to validate something and keep the modal open even if it
    // requested to be closed
    this.closeModal();
  }

  handleFeedbackButton() {
    this.setState({
      modalIsOpen: true,
    });
  }

  renderModal() {
    var modalStyle = {
      content : {
        display               : 'block',
        overflow              : 'auto',
        top                   : '40%',
        left                  : '50%',
        maxWidth              : '60%',
        transform             : 'translate(-50%, -50%)',
        bottom                : 'none',
      }
    };

    var modalOutputStyles = {
      // display     :'table-row',
      position    :'absolute',
      padding     :'0',
      margin      :'0',
      top         :'60px',
    };

    // Maybe have one Modal somewhere global and send components as arguments
    return (
      <Modal
        id='nav'
        isOpen={this.state.modalIsOpen}
        onRequestClose={this.closeModal}
        style={modalStyle} >

        <div>
          <div className="modal-header">
            <button type="button" className="close" onClick={this.handleModalCloseRequest}>
              <span aria-hidden="true">&times;</span>
              <span className="sr-only">Close</span>
            </button>
            <h4 className="modal-title">{"Contact Us"}</h4>
          </div>
          <div className="modal-body">

            <div className="container" >
              <div className="row" style={{textAlign: 'left'}}>
                <input type="text" placeholder="email" />
              </div>
              <div className="row">
                <textarea></textarea>
              </div>
              <div className="row">
                <button onClick={this.handleFeedbackButton}>Send Feedback</button>
              </div>
            </div>

          </div>
        </div>
      </Modal>
    );
  }

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
        <button className="feedback-button" onClick={this.handleFeedbackButton.bind(this)}>Contact</button>
        {this.renderModal()}
      </div>
    );
  }
}

