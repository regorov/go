import numpy
import random
import pygame
import pygame.surface
import pygame.transform
import pygame.rect
import pygame.draw

pygame.init()

def draw_shape():
    shape = random.randrange(2)
    
    surf = pygame.surface.Surface((1200, 1200))
    surf.fill((255, 255, 255))
    
    if shape == 0: # Square
        side = (random.random() * 9.0) + 3.0
        offset = int((side / 2.0) * 100)
        intside = int(side * 100)
        pygame.draw.rect(surf, (0, 0, 0), pygame.rect.Rect(offset, offset, intside, intside))
    
    elif shape == 1: # Circle
        radius = ((random.random() * 9.0) + 3.0) / 2.0
        pos = (surf.get_width() / 2, surf.get_height() / 2)
        pygame.draw.circle(surf, (0, 0, 0), pos, int(radius * 100))
    
    dest = pygame.surface.Surface((12, 12))
    pygame.transform.scale(surf, (12, 12), dest)
    
    return pygame.surfarray.array2d(dest), shape

if __name__ == "__main__":
    print "package main"
    print
    
    for name in ["TrainingData", "TestData"]:
        print "var %s = []Case{" % name
        
        for i in range(10):
            array, shape = draw_shape()
            
            is_rect = "0.9" if shape == 0 else "0.1"
            is_circle = "0.9" if shape == 1 else "0.1"
            
            print "  Case{"
            print "    []float64{"
            
            for x in range(12):
                line = "      "
                for y in range(12):
                    v = array[x][y]
                    if v == 0: # black
                        line += "0.9, "
                    else:
                        line += "0.1, "
                
                print line.rstrip()
            
            print "    },"
            print "    []float64{%s, %s}," % (is_rect, is_circle)
            print "  },"
        
        print "}"
        print