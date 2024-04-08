import React from 'react';
import { Navbar, Nav, Container } from 'react-bootstrap';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import logo from '../img/logo.png';
import Create from '../pages/Create';
import Create2 from '../pages/Create2';
import CommandList from '../pages/CommandList';
import Command from '../pages/Command';
import Stop from '../pages/Stop';

const Header = () => {

    return (
        <div>
            <Navbar className="main-header" fixed="top" collapseOnSelect expand="lg" bg="light" variant="light" >
                <Container>
                    <Navbar.Brand className="home-link">
                        <img
                            src={logo}
                            height="31"
                            width="63"
                            className="d-inline-block align-top"
                            alt="logo"
                        /> Script Executor
                    </Navbar.Brand>
                    <Navbar.Toggle aria-controls="responsive-navbar-nav" />
                    <Navbar.Collapse id="responsive-navbar-nav">
                        <Nav className="me-auto">
                            <Nav.Link className="ms-3" href="/" > Create New Command </Nav.Link>
                        </Nav>
                        <Nav className="me-auto">
                            <Nav.Link className="ms-3" href="/list" > Show Commands List </Nav.Link>
                        </Nav>
                        <Nav className="me-auto">
                            <Nav.Link className="ms-3" href="/cmd" > Show One Command </Nav.Link>
                        </Nav>
                        <Nav className="me-auto">
                            <Nav.Link className="ms-3" href="/stop" > Stop Command </Nav.Link>
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>

            <Router>
                <Routes>
                    <Route path="/" element={<Create />} />
                    <Route path="/upload" element={<Create2 />} />
                    <Route path="/list" element={<CommandList />} />
                    <Route path="/cmd" element={<Command />} />
                    <Route path="/stop" element={<Stop />} />
                </Routes>
            </Router>
        </div>
    )
}

export default Header