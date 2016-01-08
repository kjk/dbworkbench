import React from 'react';
import Modal from 'react-modal';
import Output from './Output.jsx';
import * as api from './api.js';
import * as action from './action.js';

export default class DatabaseMenuDropdown extends React.Component {
  constructor(props, context) {
    super(props, context);
    this.handleClose = this.handleClose.bind(this);
    this.handleConnection = this.handleConnection.bind(this);
    this.handleActivity = this.handleActivity.bind(this);
    this.handleDisconnect = this.handleDisconnect.bind(this);
    this.handleModalCloseRequest = this.handleModalCloseRequest.bind(this);

    this.state = {
      modalIsOpen: false,
      selectedView: '',
      results: null
    };
  }

  handleClose() {
    this.setState({
      modalIsOpen: false
    });
  }

  handleModalCloseRequest() {
    // opportunity to validate something and keep the modal open even if it
    // requested to be closed
    this.handleClose();
  }

  handleConnection() {
    console.log('handleConnection');

    const connId = this.props.connectionId;
    api.getConnectionInfo(connId, (data) => {
      this.setState({
        results: data,
        modalIsOpen: true,
        selectedView: 'Connection Info'
      });
    });
  }

  handleActivity() {
    console.log('handleActivity');

    const connId = this.props.connectionId;
    api.getActivity(connId, (data) => {
      console.log('getActivity: ', data);
      this.setState({
        results: data,
        modalIsOpen: true,
        selectedView: 'Activity'
      });
    });
  }

  handleDisconnect() {
    console.log('handleDisconnect');
    action.disconnectDatabase();
  }

  render() {
    var modalStyle = {
      content: {
        display: 'block',
        overflow: 'auto',
        top: '40%',
        left: '50%',
        maxWidth: '60%',
        transform: 'translate(-50%, -50%)',
        bottom: 'none',
      }
    };

    var modalOutputStyles = {
      // display     :'table-row',
      position: 'absolute',
      padding: 0,
      margin: 0,
      top: 60,
    };


    var appElement = document.getElementById('main');
    Modal.setAppElement(appElement);

    return (
      <div id="deneme" className='dropdown-window'>
        <div className="list-group">
          <a href="#" className="list-group-item" onClick={ this.props.handleRefresh }>Refresh Tables</a>
          <a href="#" className="list-group-item" onClick={ this.handleConnection }>Connection Info</a>
          <a href="#" className="list-group-item" onClick={ this.handleActivity }>Activity</a>
          <a href="#" className="list-group-item" onClick={ this.handleDisconnect }>Disconnect</a>
        </div>
        <Modal id='nav'
          isOpen={ this.state.modalIsOpen }
          onRequestClose={ this.handleClose }
          style={ modalStyle }>
          <div>
            <div className="modal-header">
              <button type="button" className="close" onClick={ this.handleModalCloseRequest }>
                <span aria-hidden="true">&times;</span>
                <span className="sr-only">Close</span>
              </button>
              <h4 className="modal-title">{ this.state.selectedView }</h4>
            </div>
            <div className="modal-body">
              <Output style={ modalOutputStyles } results={ this.state.results } isSidebar />
            </div>
          </div>
        </Modal>
      </div>
      );
  }
}

DatabaseMenuDropdown.propTypes = {
  handleRefresh: React.PropTypes.func.isRequired,
  connectionId: React.PropTypes.number.isRequired
};
