# gpsdec

## Simulate GPS Distance Calculation 


![alt text](https://i.imgur.com/X2yUh3l.png, "Example")

### Features
- Move characters P & Q to respective locations
- Add Buildings to complicate signal strength
- Add weather effects to complicate signals
- Adjust amount of satellites or introduce gps clock drift

With the advent of the modern Global Positioning System and all its use cases,
there still exists today an overwhelming issue with the overall size of distances
calculated by the space-based radionavigation system.

The miscalculated distances are a result of measurement and interpolation error,
this error directly affects movement data.
I will be providing a simulation for the miscalculated distances to show visually
how errors manifest themselves. 

The simulation will focus on the detailed errors
claimed by Ranacher, Brunauer, Trutschnig, Van der Spek, and Siegfried
in ”Why GPS makes distances bigger than they are”.

The article in question claims several influencing factors affect GPS measurement
error.


• Propagation delay: External factors affecting GPS signal.

• Drift in the GPS clock: A di↵erence between onboard clocks of multiple
GPS satellites.

• Ephemeris error: External factors affecting the calculated orbital position
of a satellite.

• Hardware error: Issues with signal production predicated on faulty
hardware.

• Multipath propagation: Buildings and other objects affecting travel
time of GPS signals.

• Satellite Geometry: Less than valid positioning of satellites reducing
calculated positional accuracy.

Sources:
Peter Ranacher, Richard Brunauer, Wolfgang Trutschnig, Stefan Van der Spek
Siegfried Reich (2016) Why GPS makes distances bigger than they are, International
Journal of Geographical Information Science, 30:2, 316-333, DOI:
10.1080/13658816.2015.1086924
